package main

import (
	_ "embed"
	"fmt"
	ics "github.com/arran4/golang-ical"
	"github.com/google/uuid"
	"os"
	"strings"
	"time"
)

//go:embed "wpf.htm"
var wpf string

//go:embed "mm.html"
var mm string

type lecture struct {
	start    time.Time
	end      time.Time
	category string
	number   int
	short    string
	name     string
	prof     string
	location string
}

func doStuff(all string) []lecture {
	lines := strings.Split(all, "\n")

	lectures := make([]lecture, len(lines))

	for i, line := range lines {
		if line == "" {
			continue
		}
		fields := strings.Split(line, "<br>")
		var k string
		var n int
		var d, m, y int
		var h1, m1, h2, m2 int

		fmt.Sscanf(fields[0], "%s %d.%d.%d", &k, &d, &m, &y)
		fmt.Sscanf(fields[1], "%d:%d - %d:%d Uhr", &h1, &m1, &h2, &m2)
		fmt.Sscanf(fields[3], "%d", &n)

		st := time.Date(y, time.Month(m), d, h1, m1, 0, 0, time.Local)
		en := time.Date(y, time.Month(m), d, h2, m2, 0, 0, time.Local)

		lectures[i] = lecture{
			start:    st,
			end:      en,
			category: fields[2],
			number:   n,
			short:    fields[4],
			name:     fields[5],
			prof:     fields[6],
			location: fields[7],
		}
	}
	return lectures
}

func main() {
	l1 := doStuff(wpf)
	l2 := doStuff(mm)

	calendar := ics.NewCalendar()

	lectures := append(l1, l2...)

	println(fmt.Sprintf("Total: %d, EC: %d + Mandatory: %d", len(lectures), len(l1), len(l2)))

	for _, lec := range lectures {

		desc := fmt.Sprintf("Veranstaltung: %s\nStudiengruppe: %s\nVeranstaltung: %s\nV.-Nummer 1: %d\nAnmerkung: -\nStatus: -\nDozent: %s\nRaum: %s", lec.name, lec.category, lec.short, lec.number, lec.prof, lec.location)

		event := calendar.AddEvent(uuid.NewString())
		event.SetSummary(lec.name + "/" + lec.category)
		event.SetDescription(desc)
		event.SetLocation(lec.location)
		event.SetStartAt(lec.start)
		event.SetEndAt(lec.end)
	}

	f, _ := os.Create("ws.ics")

	calendar.SerializeTo(f)
}
