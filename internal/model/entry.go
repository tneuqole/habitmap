package model

type Entry struct {
	ID      int        `json:"id"`
	HabitID int        `json:"habitId"`
	Date    CustomDate `json:"date" binding:"required"`
}
