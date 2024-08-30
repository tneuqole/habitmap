package model

type Entry struct {
	ID      int        `json:"id"`
	HabitID int        `json:"habitId"`
	Date    CustomDate `json:"date" binding:"required"`
}

type EntryDateRangeQuery struct {
	HabitID   int    `query:"habitId"`
	StartDate string `query:"startDate"`
	EndDate   string `query:"endDate"`
}
