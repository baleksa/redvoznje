package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly/v2"
)

const (
	userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36"
)

var allowedDomains = []string{"busevi.com", "www.busevi.com"}

type cityId string

const (
	BELGRADEID cityId = "bg"
	NISID      cityId = "ni"
)

type publicTransportKind string

const (
	BUS     publicTransportKind = "bus"
	TRAM    publicTransportKind = "tram"
	TROLLEY publicTransportKind = "trolley"
	MINIBUS publicTransportKind = "minibus"
	ECOBUS  publicTransportKind = "ecobus"
)

// var cityIds = []cityId{BELGRADEID, NISID}

type city struct {
	id   cityId
	name string
	t    []publicTransport
}

type publicTransport struct {
	kind    publicTransportKind
	postReq postReq
}

type postReq struct {
	url               string
	postDataScrapeUrl string
}

var bg city = city{BELGRADEID, "Beograd", []publicTransport{
	{BUS, postReq{
		"https://www.busevi.com/wp-admin/admin-ajax.php",
		"https://www.busevi.com/gradski-prevoz-beograd-autobuske-linije-brojevi-linija/"}},
	{TRAM, postReq{
		"https://www.busevi.com/wp-admin/admin-ajax.php",
		"https://www.busevi.com/gradski-prevoz-beograd-tramvajske-linije-brojevi-linija/"}},
	{ECOBUS, postReq{
		"https://www.busevi.com/wp-admin/admin-ajax.php",
		"https://www.busevi.com/gradski-prevoz-beograd-elektro-bus-linija-brojevi-linija/"}},
	{MINIBUS, postReq{
		"https://www.busevi.com/wp-admin/admin-ajax.php",
		"https://www.busevi.com/gradski-prevoz-beograd-minibus-linije-brojevi-linija/"}},
	{TROLLEY, postReq{
		"https://www.busevi.com/wp-admin/admin-ajax.php",
		"https://www.busevi.com/gradski-prevoz-beograd-trolejbuske-linije-brojevi-linija/"}}}}

func (pr postReq) PostData(c *colly.Collector) map[string]string {
	const (
		divLinesContainerCSSSel = "div.vc_grid-container.vc_clearfix.wpb_content_element.vc_basic_grid.center"
	)

	postData := map[string]string{}
	postDataJson := struct {
		Page_id        int    `json:"page_id,omitempty"`
		Style          string `json:"style,omitempty"`
		Action         string `json:"action,omitempty"`
		Shortcode_id   string `json:"shortcode_id,omitempty"`
		Items_per_page string `json:"items_per_page,omitempty"`
		Tag            string `json:"tag,omitempty"`
	}{}

	found := false

	c.OnHTML(divLinesContainerCSSSel, func(h *colly.HTMLElement) {
		if found {
			return
		}
		found = true
		postDataJSONString := h.Attr("data-vc-grid-settings")
		if err := json.Unmarshal([]byte(postDataJSONString), &postDataJson); err != nil {
			log.Fatalf("Error unmarshaling post data: %v\n", err)
		}
	})

	if err := c.Visit(pr.postDataScrapeUrl); err != nil {
		log.Fatalf("Error scraping %s for post data: %v\n", pr.postDataScrapeUrl, err)
	}
	c.Wait()

	postData["action"] = postDataJson.Action
	postData["vc_action"] = postDataJson.Action
	postData["tag"] = postDataJson.Tag
	postData["data[visible_pages]"] = "5"
	postData["data[page_id]"] = fmt.Sprint(postDataJson.Page_id)
	postData["data[style]"] = postDataJson.Style
	postData["data[action]"] = postDataJson.Action
	postData["data[shortcode_id]"] = postDataJson.Shortcode_id
	postData["data[items_per_page]"] = postDataJson.Items_per_page
	postData["data[tag]"] = postDataJson.Tag
	postData["vc_post_id"] = fmt.Sprint(postDataJson.Page_id)

	return postData
}

//	var bg city = city{BELGRADEID, "Beograd", []publicTransport{
//		{BUS, postReq{
//			"https://www.busevi.com/wp-admin/admin-ajax.php",
//			map[string]string{
//				"action":               "vc_get_vc_grid_data",
//				"vc_action":            "vc_get_vc_grid_data",
//				"tag":                  "vc_basic_grid",
//				"data[visible_pages]":  "5",
//				"data[page_id]":        "67518",
//				"data[style]":          "lazy",
//				"data[action]":         "vc_get_vc_grid_data",
//				"data[shortcode_id]":   "1700841575356-2d8044ea-125b-5",
//				"data[items_per_page]": "6",
//				"data[tag]":            "vc_basic_grid",
//				"vc_post_id":           "67518"}}},
//		{TRAM, postReq{
//			"https://www.busevi.com/wp-admin/admin-ajax.php",
//			map[string]string{
//				"action":              "vc_get_vc_grid_data",
//				"vc_action":           "vc_get_vc_grid_data",
//				"tag":                 "vc_basic_grid",
//				"data[visible_pages]": "5",
//				"data[page_id]":       "65711",
//				"data[style]":         "all",
//				"data[action]":        "vc_get_vc_grid_data",
//				"data[shortcode_id]":  "1700841327151-43306baf-a29d-3",
//				"data[tag]":           "vc_basic_grid",
//				"vc_post_id":          "65711"}}},
//		{ECOBUS, postReq{
//			"https://www.busevi.com/wp-admin/admin-ajax.php",
//			map[string]string{
//				"action":              "vc_get_vc_grid_data",
//				"vc_action":           "vc_get_vc_grid_data",
//				"tag":                 "vc_basic_grid",
//				"data[visible_pages]": "5",
//				"data[page_id]":       "65722",
//				"data[style]":         "all",
//				"data[action]":        "vc_get_vc_grid_data",
//				"data[shortcode_id]":  "1700841380352-15c43a0b-6fe5-5",
//				"data[tag]":           "vc_basic_grid",
//				"vc_post_id":          "65722"}}},
//		{MINIBUS, postReq{
//			"https://www.busevi.com/wp-admin/admin-ajax.php",
//			map[string]string{
//				"action":              "vc_get_vc_grid_data",
//				"vc_action":           "vc_get_vc_grid_data",
//				"tag":                 "vc_basic_grid",
//				"data[visible_pages]": "5",
//				"data[page_id]":       "65728",
//				"data[style]":         "all",
//				"data[action]":        "vc_get_vc_grid_data",
//				"data[shortcode_id]":  "1700841409243-ebae0051-9b12-9",
//				"data[tag]":           "vc_basic_grid",
//				"vc_post_id":          "65728"}}},
//		{TROLLEY, postReq{
//			"https://www.busevi.com/wp-admin/admin-ajax.php",
//			map[string]string{
//				"action":              "vc_get_vc_grid_data",
//				"vc_action":           "vc_get_vc_grid_data",
//				"tag":                 "vc_basic_grid",
//				"data[visible_pages]": "5",
//				"data[page_id]":       "65715",
//				"data[style]":         "all",
//				"data[action]":        "vc_get_vc_grid_data",
//				"data[shortcode_id]":  "1700841359426-29bdc697-db18-7",
//				"data[tag]":           "vc_basic_grid",
//				"vc_post_id":          "65715"}}}}}

// var ni city = city{NISID, "Niš", []publicTransport{}}

var cities = []*city{&bg}

func scrapeTransportLinesLinks() map[string]string {
	const (
		transportLineContainerCSSSel = "div.vc_gitem-zone.vc_gitem-zone-a.linija"
	)

	transportLines := make(map[string]string)

	c := colly.NewCollector(colly.AllowedDomains(allowedDomains...), colly.Async(true))
	c.UserAgent = userAgent

	vnconce, err := getVnconce("https://www.busevi.com", c.Clone())
	if err != nil {
		log.Fatal(err)
	}

	c.OnHTML(transportLineContainerCSSSel, func(h *colly.HTMLElement) {
		lineLink := h.ChildAttr("a", "href")
		lineId := strings.TrimSpace(h.Text)
		transportLines[lineId] = lineLink
	})

	for _, city := range cities {
		for _, t := range city.t {
			log.Printf("Scraping %s lines in %s.\n", t.kind, city.name)

			postData := t.postReq.PostData(c.Clone())
			postData["_vcnonce"] = vnconce

			log.Printf("Post data map: %v\n", postData)

			if err := c.Post(t.postReq.url, postData); err != nil {
				log.Fatalf("Colly post request to %s with data %vf failed: %v\n", t.postReq.url, postData, err)
			}
			c.Wait()
		}
	}

	log.Println("Got all links!")

	return transportLines
}

func scrapeLine(url string) (*TransportLine, error) {
	log.Printf("Scraping line schedule from url: %s\n", url)

	const (
		LineIdSpanCssSel              = "table.alignleft:nth-child(2) > tbody:nth-child(1) > tr:nth-child(1) > td:nth-child(1) > span:nth-child(1)"
		LinePlacesTableCssSel         = "table.alignleft:nth-child(3)"
		LinePlacesCssSel              = "table.alignleft:nth-child(3) > tbody:nth-child(1) > tr:nth-child(1) > td:nth-child(1) > div:nth-child(1) > span:nth-child(1)"
		TimetableCSSSel               = "table.tablepress-initially-hidden tbody"
		TimetableInitiallyShownCSSSel = "table.tablepress tbody"
	)

	b := TransportLine{}
	c := colly.NewCollector(colly.AllowedDomains(allowedDomains...))

	timetables := []timetable{}
	timetablesInitiallyShown := []timetable{}
	stops := [][]stop{}
	tags := []string{}

	c.UserAgent = userAgent

	c.OnHTML(LineIdSpanCssSel, func(h *colly.HTMLElement) {
		b.Id = strings.TrimSpace(h.Text)
	})

	c.OnHTML(LinePlacesTableCssSel, func(h *colly.HTMLElement) {
		if b.Places != nil {
			return
		}
		b.Places = h.ChildTexts("span")
		for i, place := range b.Places {
			b.Places[i] = strings.TrimSpace(strings.Trim(place, "›‹·"))
		}
	})
	c.OnHTML(TimetableCSSSel, func(h *colly.HTMLElement) {
		tt := timetable{}
		h.ForEach("tr", func(i int, h *colly.HTMLElement) {
			const nClmns = 4
			tr := h.ChildTexts("td")
			if len(tr) != nClmns {
				// Table row should have 4 clmns hour, workday, saturday, sunday.
				return
			}
			tt = append(tt, parseTimelineRowFromTableRow(tr))
		})
		timetables = append(timetables, tt)
	})

	c.OnHTML(TimetableInitiallyShownCSSSel, func(h *colly.HTMLElement) {
		tt := timetable{}
		h.ForEach("tr", func(i int, h *colly.HTMLElement) {
			const nclmns = 4
			tr := h.ChildTexts("td")
			if len(tr) != nclmns {
				// Table row should have 4 clmns hour, workday, saturday, sunday.
				return
			}
			tt = append(tt, parseTimelineRowFromTableRow(tr))
		})
		timetablesInitiallyShown = append(timetablesInitiallyShown, tt)
	})

	c.OnHTML("ul", func(h *colly.HTMLElement) {
		s := []stop{}
		if !strings.Contains(h.Text, "min./") {
			return
		}
		h.ForEach("li", func(i int, h *colly.HTMLElement) {
			const nameStartInd = 3
			x := stop{}

			fs := strings.Fields(h.Text)
			x.Name = strings.Join(fs[nameStartInd:], " ")
			x.Id = h.ChildText("a")
			x.Zone = h.ChildAttr("img", "title")

			s = append(s, x)
		})
		stops = append(stops, s)
	})

	c.OnHTML("li.vc_tta-tab > a > span", func(h *colly.HTMLElement) {
		tags = append(tags, strings.TrimSpace(h.Text))
	})

	if err := c.Visit(url); err != nil {
		return nil, err
	}

	if len(timetables) == 0 {
		timetables = timetablesInitiallyShown
	}
	b.Routes = make([]route, len(timetables))
	for i := range b.Routes {
		b.Routes[i].Timetable = timetables[i]
		if i < len(stops) {
			b.Routes[i].Stops = stops[i]
		}
		if i < len(tags) {
			b.Routes[i].Tag = tags[i]
		}
	}

	return &b, nil
}

func parseTimelineRowFromTableRow(tr []string) row {
	// len(tr) is 4
	r := row{}
	fmt.Sscanf(tr[0], "%d", &r.H)
	r.Wd = parseMinutes(tr[1])
	r.Sat = parseMinutes(tr[2])
	r.Sun = parseMinutes(tr[3])
	return r
}

func parseMinutes(td string) []min {
	fs := strings.Fields(td)
	if len(fs) == 0 || fs[0] == "-" {
		return nil
	}
	ms := make([]min, len(fs))
	for i, f := range fs {
		fmt.Sscanf(f, "%d", &ms[i])
	}
	return ms
}

func getVnconce(url string, c *colly.Collector) (string, error) {
	const (
		vncOnceCSSSel = ".vc_custom_1700841099371"
	)

	var vnconce string

	c.OnHTML(vncOnceCSSSel, func(h *colly.HTMLElement) {
		s := h.Attr("data-vc-public-nonce")
		if s != "" {
			vnconce = s
		}
	})

	if err := c.Visit(url); err != nil {
		return "", err
	}

	c.Wait()

	if vnconce == "" {
		return "", fmt.Errorf("vcnonce not found.")
	}

	log.Printf("vnconce=%s\n", vnconce)

	return vnconce, nil
}
