package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tneuqole/habitmap/internal/model"
	"github.com/tneuqole/habitmap/internal/templates"
)

var validate = NewValidate()

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
	return Render(c, templates.Habits(habits))
}

func (h *HabitHandler) GetNewHabitForm(c echo.Context) error {
	return Render(c, templates.NewHabit(templates.NewHabitFormData{}))
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

	return Render(c, templates.HabitView(habit))
}

type PostHabitForm struct {
	Name string `form:"name" validate:"required,notblank,min=1,max=32"`
}

func (h *HabitHandler) PostHabit(c echo.Context) error {
	form := PostHabitForm{}
	if err := c.Bind(&form); err != nil {
		return err
	}

	err := validate.Struct(&form)
	if err != nil {
		errors := ParseValidationErrors(err)
		data := templates.NewHabitFormData{
			Name:   form.Name,
			Errors: errors,
		}
		return Render(c, templates.NewHabit(data))
	}

	habit, err := h.queries.CreateHabit(c.Request().Context(), form.Name)
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/habits/%d", habit.ID))
}
