package templates

import (
	"time"
)

type HabitFormData struct {
	Name   string
	ID     int64
	Errors map[string]string
}

func UpdateDate(date, view string, inc int) string {
	t, err := time.Parse("2006-01", date)
	if err != nil {
		// TODO log error
		return ""
	}

	switch view {
	case "year":
		t = t.AddDate(inc, 0, 0)
	case "month":
		t = t.AddDate(0, inc, 0)
	}

	return t.Format("2006-01")
}

func FormatDate(date, view string) string {
	t, err := time.Parse("2006-01", date)
	if err != nil {
		// TODO log error
		return ""
	}

	switch view {
	case "year":
		return t.Format("2006")
	case "month":
		return t.Format("January 2006")
	default:
		return ""
	}
}
