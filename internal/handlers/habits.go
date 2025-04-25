package handlers

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/labstack/echo/v4"
	"github.com/tneuqole/habitmap/internal/model"
	"github.com/tneuqole/habitmap/internal/templates"
	"github.com/tneuqole/habitmap/internal/templates/forms"
	"github.com/tneuqole/habitmap/internal/templates/pages"
)

// HabitHandler handles HTTP requests related to habits
type HabitHandler struct {
	*BaseHandler
}

// NewHabitHandler creates a new HabitHandler
func NewHabitHandler(bh *BaseHandler) *HabitHandler {
	return &HabitHandler{
		BaseHandler: bh,
	}
}

// GetHabits returns all habits from the database
func (h *HabitHandler) GetHabits(c echo.Context) error {
	habits, err := h.Queries.GetHabits(c.Request().Context())
	if err != nil {
		return err
	}

	return h.render(c, pages.Habits(habits))
}

type getHabitParams struct {
	HabitID int64 `param:"id"`
}

// GetHabit returns a single habit by id
func (h *HabitHandler) GetHabit(c echo.Context) error {
	params := getHabitParams{}
	if err := c.Bind(&params); err != nil {
		return err
	}

	habit, err := h.Queries.GetHabit(c.Request().Context(), params.HabitID)
	if err != nil {
		return err
	}

	entries, err := h.Queries.GetEntriesForHabit(c.Request().Context(), habit.ID)
	if err != nil {
		return err
	}

	// group entries by month
	entriesMonthMap := make(map[string][]model.Entry)
	for _, entry := range entries {
		key := entry.EntryDate[:7] // YYYY-MM
		entriesMonthMap[key] = append(entriesMonthMap[key], entry)
	}

	// sort months from past to present
	var sortedMonths []string
	for key := range entriesMonthMap {
		sortedMonths = append(sortedMonths, key)
	}
	sort.Strings(sortedMonths)

	// parse entries into 2D arrays representing a month
	entriesForMonths := make(map[string][][]model.Entry)
	for monthStr, entries := range entriesMonthMap {
		entriesForMonths[monthStr] = h.generateMonth(monthStr, entries)
	}
	return h.render(c, pages.Habit(habit, sortedMonths, entriesForMonths))
}

// DeleteHabit deletes a single habit by id
func (h *HabitHandler) DeleteHabit(c echo.Context) error {
	params := getHabitParams{}
	if err := c.Bind(&params); err != nil {
		return err
	}

	err := h.Queries.DeleteHabit(c.Request().Context(), params.HabitID)
	if err != nil {
		return err
	}

	c.Response().Header().Add("Hx-Redirect", "/habits")

	return nil
}

// GetCreateHabitForm renders a form for creating a new habit
func (h *HabitHandler) GetCreateHabitForm(c echo.Context) error {
	return h.render(c, forms.CreateHabit(templates.HabitFormData{}))
}

type createHabitForm struct {
	Name string `form:"name" validate:"required,notblank,min=1,max=32"`
}

// PostHabit processes a form for creating a new habit
func (h *HabitHandler) PostHabit(c echo.Context) error {
	form := createHabitForm{}
	if err := c.Bind(&form); err != nil {
		return err
	}

	err := validate.Struct(&form)
	if err != nil {
		errors := parseValidationErrors(err)
		data := templates.HabitFormData{
			Name:   form.Name,
			Errors: errors,
		}
		return h.render(c, forms.CreateHabit(data))
	}

	habit, err := h.Queries.CreateHabit(c.Request().Context(), form.Name)
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/habits/%d", habit.ID))
}

// GetUpdateHabitForm renders a form for updating a habit
func (h *HabitHandler) GetUpdateHabitForm(c echo.Context) error {
	params := getHabitParams{}
	if err := c.Bind(&params); err != nil {
		return err
	}
	return h.render(c, forms.UpdateHabit(templates.HabitFormData{ID: params.HabitID}))
}

type updateHabitForm struct {
	HabitID int64 `param:"id"`
	createHabitForm
}

// PostUpdateHabit processes a form for updating a habit
func (h *HabitHandler) PostUpdateHabit(c echo.Context) error {
	form := updateHabitForm{}
	if err := c.Bind(&form); err != nil {
		return err
	}

	err := validate.Struct(&form)
	if err != nil {
		errors := parseValidationErrors(err)
		data := templates.HabitFormData{
			Name:   form.Name,
			Errors: errors,
		}
		return h.render(c, forms.UpdateHabit(data))
	}

	habit, err := h.Queries.UpdateHabit(c.Request().Context(), model.UpdateHabitParams{
		Name: form.Name,
		ID:   form.HabitID,
	})
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/habits/%d", habit.ID))
}
