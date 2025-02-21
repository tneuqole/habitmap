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

type HabitHandler struct {
	queries *model.Queries
}

func NewHabitHandler(queries *model.Queries) *HabitHandler {
	return &HabitHandler{
		queries: queries,
	}
}

func (h *HabitHandler) GetHabits(c echo.Context) error {
	habits, err := h.queries.GetHabits(c.Request().Context())
	if err != nil {
		return err
	}

	c.Logger().Info(habits)
	return Render(c, pages.Habits(habits))
}

type GetHabitParams struct {
	HabitID int64 `param:"id"`
}

func (h *HabitHandler) GetHabit(c echo.Context) error {
	params := GetHabitParams{}
	if err := c.Bind(&params); err != nil {
		return err
	}

	habit, err := h.queries.GetHabit(c.Request().Context(), params.HabitID)
	if err != nil {
		return err
	}

	entries, err := h.queries.GetEntriesForHabit(c.Request().Context(), habit.ID)
	if err != nil {
		return err
	}

	months := make(map[string][]model.Entry)
	for _, entry := range entries {
		key := entry.EntryDate[:7] // YYYY-MM
		months[key] = append(months[key], entry)
	}

	var keys []string
	for key := range months {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	return Render(c, pages.Habit(habit, keys, months))
}

func (h *HabitHandler) DeleteHabit(c echo.Context) error {
	params := GetHabitParams{}
	if err := c.Bind(&params); err != nil {
		return err
	}

	err := h.queries.DeleteHabit(c.Request().Context(), params.HabitID)
	if err != nil {
		return err
	}

	c.Response().Header().Add("Hx-Redirect", "/habits")

	return nil
}

func (h *HabitHandler) GetCreateHabitForm(c echo.Context) error {
	return Render(c, forms.CreateHabit(templates.HabitFormData{}))
}

type CreateHabitForm struct {
	Name string `form:"name" validate:"required,notblank,min=1,max=32"`
}

func (h *HabitHandler) PostHabit(c echo.Context) error {
	form := CreateHabitForm{}
	if err := c.Bind(&form); err != nil {
		return err
	}

	err := validate.Struct(&form)
	if err != nil {
		errors := ParseValidationErrors(err)
		data := templates.HabitFormData{
			Name:   form.Name,
			Errors: errors,
		}
		return Render(c, forms.CreateHabit(data))
	}

	habit, err := h.queries.CreateHabit(c.Request().Context(), form.Name)
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/habits/%d", habit.ID))
}

func (h *HabitHandler) GetUpdateHabitForm(c echo.Context) error {
	params := GetHabitParams{}
	if err := c.Bind(&params); err != nil {
		return err
	}
	return Render(c, forms.UpdateHabit(templates.HabitFormData{ID: params.HabitID}))
}

type UpdateHabitForm struct {
	HabitID int64 `param:"id"`
	CreateHabitForm
}

func (h *HabitHandler) PostUpdateHabit(c echo.Context) error {
	form := UpdateHabitForm{}
	if err := c.Bind(&form); err != nil {
		return err
	}

	err := validate.Struct(&form)
	if err != nil {
		errors := ParseValidationErrors(err)
		data := templates.HabitFormData{
			Name:   form.Name,
			Errors: errors,
		}
		return Render(c, forms.UpdateHabit(data))
	}

	habit, err := h.queries.UpdateHabit(c.Request().Context(), model.UpdateHabitParams{
		Name: form.Name,
		ID:   form.HabitID,
	})
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/habits/%d", habit.ID))
}
