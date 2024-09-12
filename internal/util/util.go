package util

import (
	"time"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/tneuqole/habitmap/internal/model"
)

func Render(c echo.Context, component templ.Component) error {
	return component.Render(c.Request().Context(), c.Response())
}

// TODO entries should be sorted by date
func GenerateMonth(t time.Time, entries []model.Entry) [][]model.Entry {
	var month [][]model.Entry
	var week = make([]model.Entry, 7)

	date := t.AddDate(0, 0, -t.Day()+1)
	daysInMonth := date.AddDate(0, 1, -1).Day()

	habitID := entries[0].HabitID
	entryIdx := 0
	dayOfWeek := int(date.Weekday())
	for day := date.Day(); day <= daysInMonth; {
		for ; dayOfWeek < 7 && day <= daysInMonth; dayOfWeek++ {
			if entryIdx < len(entries) && entries[entryIdx].Date.Time() == date {
				week[dayOfWeek] = entries[entryIdx]
				entryIdx++
			} else {
				entry := model.NewEntry()
				entry.HabitID = habitID
				entry.Date = model.CustomDate(date)
				week[dayOfWeek] = entry
			}
			date = date.AddDate(0, 0, 1)
			day++
		}
		month = append(month, week)
		week = make([]model.Entry, 7)
		dayOfWeek = 0
	}

	return month
}
