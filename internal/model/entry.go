package model

type Entry struct {
	ID      int        `form:"id" json:"id"`
	HabitID int        `form:"habitId" json:"habitId"`
	Date    CustomDate `form:"date" json:"date" binding:"required"`
}

func NewEntry() Entry {
	return Entry{
		ID: -1,
	}
}
