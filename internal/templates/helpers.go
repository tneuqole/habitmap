package templates

import (
	"fmt"
	"time"

	"github.com/tneuqole/habitmap/internal/model"
)

// m = "YYYY-DD"
func GenerateMonth(monthStr string, entries []model.Entry) [][]model.Entry {
	var month [][]model.Entry
	week := make([]model.Entry, 7)

	t, err := time.Parse("2006-01", monthStr)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return month
	}

	date := t.AddDate(0, 0, -t.Day()+1)
	daysInMonth := date.AddDate(0, 1, -1).Day()

	habitID := entries[0].HabitID
	entryIdx := 0
	dayOfWeek := int(date.Weekday())
	for day := date.Day(); day <= daysInMonth; {
		for ; dayOfWeek < 7 && day <= daysInMonth; dayOfWeek++ {
			if entryIdx < len(entries) && entries[entryIdx].EntryDate == date.Unix() {
				week[dayOfWeek] = entries[entryIdx]
				entryIdx++
			} else {
				entry := model.Entry{
					HabitID:   habitID,
					EntryDate: date.Unix(),
				}
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
