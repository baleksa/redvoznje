package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly/v2"
)

const (
	vncOnceCssSel = ".vc_custom_1700841099371"

	buslineLinkCssSel      = "a.vc-zone-link"
	buslineIdCssSel        = "dim .vc_gitem-acf"
	buslineContainerCssSel = "div.vc_gitem-zone.vc_gitem-zone-a.linija"

	BusIdSpanCssSel      = "table.alignleft:nth-child(2) > tbody:nth-child(1) > tr:nth-child(1) > td:nth-child(1) > span:nth-child(1)"
	BusPlacesTableCssSel = "table.alignleft:nth-child(3)"
	BusPlacesCssSel      = "table.alignleft:nth-child(3) > tbody:nth-child(1) > tr:nth-child(1) > td:nth-child(1) > div:nth-child(1) > span:nth-child(1)"

	TimetableCSSSel         = "table.tablepress-initially-hidden tbody"
	TimetableIniShownCSSSel = "table.tablepress tbody"
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

var cityIds = []cityId{BELGRADEID, NISID}

type city struct {
	id   cityId
	name string
	t    []publicTransport
}

type publicTransport struct {
	kind        publicTransportKind
	dataPostReq postReq
}

type postReq struct {
	url         string
	postReqData map[string]string
}

var bg city = city{BELGRADEID, "Beograd", []publicTransport{
	{BUS, postReq{
		"https://www.busevi.com/wp-admin/admin-ajax.php",
		map[string]string{
			"action":               "vc_get_vc_grid_data",
			"vc_action":            "vc_get_vc_grid_data",
			"tag":                  "vc_basic_grid",
			"data[visible_pages]":  "5",
			"data[page_id]":        "67518",
			"data[style]":          "lazy",
			"data[action]":         "vc_get_vc_grid_data",
			"data[shortcode_id]":   "1700841575356-2d8044ea-125b-5",
			"data[items_per_page]": "6",
			"data[tag]":            "vc_basic_grid",
			"vc_post_id":           "67518"}}},
	{TRAM, postReq{
		"https://www.busevi.com/wp-admin/admin-ajax.php",
		map[string]string{
			"action":              "vc_get_vc_grid_data",
			"vc_action":           "vc_get_vc_grid_data",
			"tag":                 "vc_basic_grid",
			"data[visible_pages]": "5",
			"data[page_id]":       "65711",
			"data[style]":         "all",
			"data[action]":        "vc_get_vc_grid_data",
			"data[shortcode_id]":  "1700841327151-43306baf-a29d-3",
			"data[tag]":           "vc_basic_grid",
			"vc_post_id":          "65711"}}},
	{ECOBUS, postReq{
		"https://www.busevi.com/wp-admin/admin-ajax.php",
		map[string]string{
			"action":              "vc_get_vc_grid_data",
			"vc_action":           "vc_get_vc_grid_data",
			"tag":                 "vc_basic_grid",
			"data[visible_pages]": "5",
			"data[page_id]":       "65722",
			"data[style]":         "all",
			"data[action]":        "vc_get_vc_grid_data",
			"data[shortcode_id]":  "1700841380352-15c43a0b-6fe5-5",
			"data[tag]":           "vc_basic_grid",
			"vc_post_id":          "65722"}}},
	{MINIBUS, postReq{
		"https://www.busevi.com/wp-admin/admin-ajax.php",
		map[string]string{
			"action":              "vc_get_vc_grid_data",
			"vc_action":           "vc_get_vc_grid_data",
			"tag":                 "vc_basic_grid",
			"data[visible_pages]": "5",
			"data[page_id]":       "65728",
			"data[style]":         "all",
			"data[action]":        "vc_get_vc_grid_data",
			"data[shortcode_id]":  "1700841409243-ebae0051-9b12-9",
			"data[tag]":           "vc_basic_grid",
			"vc_post_id":          "65728"}}},
	{TROLLEY, postReq{
		"https://www.busevi.com/wp-admin/admin-ajax.php",
		map[string]string{
			"action":              "vc_get_vc_grid_data",
			"vc_action":           "vc_get_vc_grid_data",
			"tag":                 "vc_basic_grid",
			"data[visible_pages]": "5",
			"data[page_id]":       "65715",
			"data[style]":         "all",
			"data[action]":        "vc_get_vc_grid_data",
			"data[shortcode_id]":  "1700841359426-29bdc697-db18-7",
			"data[tag]":           "vc_basic_grid",
			"vc_post_id":          "65715"}}}}}

var ni city = city{NISID, "Niš", []publicTransport{}}

var cities = []*city{&bg}

func scrapeTransportLinesLinks() map[string]string {
	transportLines := make(map[string]string)

	c := colly.NewCollector(colly.AllowedDomains(allowedDomains...), colly.Async(true))
	c.UserAgent = userAgent

	vnconce, err := getVnconce("https://www.busevi.com", vncOnceCssSel, c.Clone())
	if err != nil {
		log.Fatal(err)
	}

	c.OnHTML(buslineContainerCssSel, func(h *colly.HTMLElement) {
		lineLink := h.ChildAttr("a", "href")
		lineId := strings.TrimSpace(h.Text)
		transportLines[lineId] = lineLink
	})

	for _, city := range cities {
		for _, t := range city.t {
			log.Printf("Scraping %s lines in %s.\n", t.kind, city.name)
			t.dataPostReq.postReqData["_vcnonce"] = vnconce
			err := c.Post(t.dataPostReq.url, t.dataPostReq.postReqData)
			if err != nil {
				log.Fatal(err)
			}
			c.Wait()
		}
	}

	log.Println("Got all links!")

	return transportLines
}

func scrapeBusLine(url string) (*BusLineSchedule, error) {
	log.Printf("Scraping line schedule from url: %s\n", url)
	b := BusLineSchedule{}
	c := colly.NewCollector(colly.AllowedDomains(allowedDomains...))

	timetables := []timetable{}
	timetablesInitiallyShown := []timetable{}
	stops := [][]stop{}
	tags := []string{}

	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36"

	c.OnHTML(BusIdSpanCssSel, func(h *colly.HTMLElement) {
		b.Id = strings.TrimSpace(h.Text)
	})

	c.OnHTML(BusPlacesTableCssSel, func(h *colly.HTMLElement) {
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
			const nclmns = 4
			tr := h.ChildTexts("td")
			if len(tr) != nclmns {
				// Table row should have 4 clmns hour, workday, saturday, sunday.
				return
			}
			tt = append(tt, parseTimelineRowFromTableRow(tr))
		})
		timetables = append(timetables, tt)
	})

	c.OnHTML(TimetableIniShownCSSSel, func(h *colly.HTMLElement) {
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
		b.Routes[i].Stops = stops[i]
		b.Routes[i].Tag = tags[i]
	}

	return &b, nil
}

func parseTimelineRowFromTableRow(tr []string) row {
	r := row{}
	fmt.Sscanf(tr[0], "%d", &r.H)
	r.Wd = parseMinutes(tr[1])
	r.Sat = parseMinutes(tr[2])
	r.Sun = parseMinutes(tr[3])
	return r
}

func parseMinutes(td string) []min {
	fs := strings.Fields(td)
	if fs[0] == "-" {
		return nil
	}
	ms := make([]min, len(fs))
	for i, f := range fs {
		fmt.Sscanf(f, "%d", &ms[i])
	}
	return ms
}

func getVnconce(url string, CSSSel string, c *colly.Collector) (string, error) {
	var vnconce string

	c.OnHTML(CSSSel, func(h *colly.HTMLElement) {
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
