package templates

import (
	"fmt"
	"time"

	"github.com/tneuqole/habitmap/internal/model"
)

// monthStr = "YYYY-DD"
func GenerateMonth(monthStr string, entries []model.Entry) [][]model.Entry {
	var month [][]model.Entry
	week := make([]model.Entry, 7)

	date, err := time.Parse("2006-01", monthStr)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return month
	}

	daysInMonth := date.AddDate(0, 1, -1).Day()

	habitID := entries[0].HabitID
	entryIdx := 0
	dayOfWeek := int(date.Weekday())
	for day := date.Day(); day <= daysInMonth; {
		for ; dayOfWeek < 7 && day <= daysInMonth; dayOfWeek++ {
			fmt.Println(day, date, date.Unix())
			if entryIdx < len(entries) && entries[entryIdx].EntryDate == date.Format("2006-01-02") {
				week[dayOfWeek] = entries[entryIdx]
				entryIdx++
			} else {
				entry := model.Entry{
					HabitID:   habitID,
					EntryDate: date.Format("2006-01-02"),
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

	fmt.Println(len(month))
	for len(month) < 6 {
		week = make([]model.Entry, 7)
		month = append(month, week)
	}
	fmt.Println(len(month))

	return month
}
