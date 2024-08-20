package main

import (
	"embed"
	"flag"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"sort"
	"strings"
)

//go:embed templates static
var tmplFs embed.FS

func getKeysFromMap(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sortStringSliceNumerically(keys)
	return keys
}

func sortStringSliceNumerically(s []string) []string {
	sort.Slice(s, func(i, j int) bool {
		var xi, xj int
		fmt.Sscanf(s[i], "%d", &xi)
		fmt.Sscanf(s[j], "%d", &xj)

		if xi != xj {
			return xi < xj
		}

		return s[i] < s[j]
	})

	return s
}

func run(args []string) error {

	publicTransportLinesLinks := scrapeTransportLinesLinks()

	lineCache := newCache()

	baseTmpl := template.Must(template.ParseFS(tmplFs, "templates/base/*"))
	templates := make(map[string]*template.Template)
	loadTemplates(tmplFs, templates, baseTmpl)

	fs := http.FileServerFS(tmplFs)
	http.Handle("/static/", fs)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(templates, w, "home", sortStringSliceNumerically(getKeysFromMap(publicTransportLinesLinks)))
	})

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Error parsing request form.", http.StatusInternalServerError)
		}
		q := strings.ToLower(r.Form.Get("q"))
		matchedLines := []string{}
		for key := range publicTransportLinesLinks {
			if strings.Contains(strings.ToLower(key), q) {
				matchedLines = append(matchedLines, key)
			}
		}
		renderTemplate(templates, w, "home", sortStringSliceNumerically(matchedLines))
	})

	http.HandleFunc("/line/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		val, ok := lineCache.get(id)
		if !ok {
			v, err := scrapeLine(publicTransportLinesLinks[id])
			if err != nil {
				log.Fatalf("Scraping error: %v\n", err)
			}
			lineCache.set(id, v)
			val = v
		}
		b := val.(*TransportLine)
		renderTemplate(templates, w, "line", b)
	})

	flagset := flag.NewFlagSet("", flag.ExitOnError)
	host := flagset.String("host", "localhost", "Address to listen on.")
	port := flagset.String("port", "8080", "Port to listen on.")
	flagset.Usage = func() {
		fmt.Printf("Usage: %s [flags]\n", args[0])
		flagset.PrintDefaults()
	}
	_ = flagset.Parse(args[1:])

	addr := net.JoinHostPort(*host, *port)
	log.Printf("Listening on %s.\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := run(os.Args); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func loadTemplates(efs embed.FS, templates map[string]*template.Template, baseTmpl *template.Template) {
	files, err := fs.Glob(efs, "templates/*.html")
	if err != nil {
		log.Fatalln("Failed to match template files.", err)
	}
	for _, file := range files {
		base := path.Base(file)
		log.Printf("Basename of a template source file => %s\n", base)
		t, err := baseTmpl.Clone()
		if err != nil {
			log.Fatalln("Failed to clone main template.", err)
		}
		templates[base], err = t.ParseFS(efs, file)
		if err != nil {
			log.Fatalf("Failed to parse template file %s. %v\n", file, err)
		}
		log.Printf("Template struct %s%s\n", templates[base].Name(), templates[base].DefinedTemplates())
	}

}

func renderTemplate(templates map[string]*template.Template, w http.ResponseWriter, name string, data any) {
	t, ok := templates[name+".html"]
	if !ok {
		log.Printf("Can't find template %s.", name)
		http.Error(w, "Internal template error.", http.StatusInternalServerError)
		return
	}
	if err := t.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
