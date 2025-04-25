package templates

// HabitFormData is used to rerender Habit forms with information
// from a previous request
type HabitFormData struct {
	Name   string
	ID     int64
	Errors map[string]string
}
