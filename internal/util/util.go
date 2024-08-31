package util

import (
	"log"
	"time"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/tneuqole/habitmap/internal/model"
)

func Render(c echo.Context, component templ.Component) error {
	return component.Render(c.Request().Context(), c.Response())
}

// TODO entries should be sorted by date
func GenerateMonth(t time.Time, entries []model.Entry) [][]string {
	log.Printf("creating month for date: %s\n", t)
	var month [][]string
	var week = make([]string, 7)

	date := t.AddDate(0, 0, -t.Day()+1)
	daysInMonth := date.AddDate(0, 1, -1).Day()

	entryIdx := 0
	dayOfWeek := int(date.Weekday())
	for day := date.Day(); day <= daysInMonth; {
		for ; dayOfWeek < 7 && day <= daysInMonth; dayOfWeek++ {
			if entryIdx < len(entries) && entries[entryIdx].Date.Time() == date {
				week[dayOfWeek] = "y"
				entryIdx++
			} else {
				week[dayOfWeek] = "n"
			}
			date = date.AddDate(0, 0, 1)
			day++
		}
		month = append(month, week)
		week = make([]string, 7)
		dayOfWeek = 0
	}

	return month
}
