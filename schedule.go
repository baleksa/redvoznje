package main

import "time"

// "time"

type TransportLine struct {
	Places []string
	Id     string
	Routes []route
}

type route struct {
	Tag       string
	Stops     []stop
	Timetable timetable
}

type stop struct {
	Name string
	Zone string
	Id   string
}

type timetable []row

type hour int

type min int

type row struct {
	H   hour
	Wd  []min // workdays
	Sat []min
	Sun []min
}

func (t timetable) FirstNextDepartures() []row {
	var firstNext []row
	now := time.Now()
	nowH := now.Hour()
	nowM := now.Minute()

	for i, r := range t {
		if r.H != hour(nowH) {
			continue
		}

		var mins []min

		switch now.Weekday() {
		case time.Saturday:
			mins = r.Sat
		case time.Sunday:
			mins = r.Sun
		default:
			mins = r.Wd
		}
		for _, m := range mins {
			if m >= min(nowM) {
				firstNext = t[i:]
			} else {
				if i < len(t)-1 {
					firstNext = t[i+1:]
				}
			}
		}
	}

	if firstNext == nil {
		firstNext = t
	}

	if len(firstNext) > 1 {
		return firstNext[:2]
	} else {
		return firstNext
	}
}
