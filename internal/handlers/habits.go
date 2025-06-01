package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/tneuqole/habitmap/internal/forms"
	"github.com/tneuqole/habitmap/internal/model"
	"github.com/tneuqole/habitmap/internal/templates/formcomponents"
	"github.com/tneuqole/habitmap/internal/templates/pages"
)

type HabitHandler struct {
	*BaseHandler
}

func NewHabitHandler(bh *BaseHandler) *HabitHandler {
	return &HabitHandler{
		BaseHandler: bh,
	}
}

func (h *HabitHandler) GetHabits(w http.ResponseWriter, r *http.Request) error {
	habits, err := h.Queries.GetHabits(r.Context())
	if err != nil {
		return h.handleDBError(err)
	}

	return h.render(w, r, pages.Habits(habits))
}

type getHabitParams struct {
	View string `form:"view" validate:"required,oneof=year month"`
	Date string `form:"date" validate:"required,yearmonth"`
}

func newGetHabitParams() getHabitParams {
	return getHabitParams{
		View: "year",
		Date: time.Now().Format("2006-01"),
	}
}

func (h *HabitHandler) GetHabit(w http.ResponseWriter, r *http.Request) error {
	habitID, err := h.getIDFromURL(r)
	if err != nil {
		return err
	}

	r, params, err := h.bindAndValidateGetHabitParams(w, r)
	if err != nil {
		return err
	}

	habit, err := h.Queries.GetHabit(r.Context(), habitID)
	if err != nil {
		return h.handleDBError(err)
	}

	entries, err := h.fetchEntriesForView(r.Context(), habitID, params.View, params.Date)
	if err != nil {
		return h.handleDBError(err)
	}

	monthKeys, err := h.generateMonths(params.View, params.Date)
	if err != nil {
		return err
	}

	months := h.groupEntriesByMonth(habitID, entries, monthKeys)

	return h.render(w, r, pages.Habit(habit, params.View, params.Date, monthKeys, months))
}

func (h *HabitHandler) DeleteHabit(w http.ResponseWriter, r *http.Request) error {
	habitID, err := h.getIDFromURL(r)
	if err != nil {
		return err
	}

	err = h.Queries.DeleteHabit(r.Context(), habitID)
	if err != nil {
		return h.handleDBError(err)
	}

	w.Header().Set("HX-Redirect", "/habits")
	w.WriteHeader(http.StatusNoContent)

	return nil
}

func (h *HabitHandler) GetCreateHabitForm(w http.ResponseWriter, r *http.Request) error {
	return h.render(w, r, formcomponents.CreateHabit(forms.CreateHabitForm{}))
}

func (h *HabitHandler) PostHabit(w http.ResponseWriter, r *http.Request) error {
	var form forms.CreateHabitForm
	if err := h.bindFormData(r, &form); err != nil {
		return err
	}

	err := validate.Struct(&form)
	if err != nil {
		form.FieldErrors = h.parseValidationErrors(err)
		return h.render(w, r, formcomponents.CreateHabit(form))
	}

	habit, err := h.Queries.CreateHabit(r.Context(), form.Name)
	if err != nil {
		return h.handleDBError(err)
	}

	http.Redirect(w, r, fmt.Sprintf("/habits/%d", habit.ID), http.StatusSeeOther)
	return nil
}

func (h *HabitHandler) GetUpdateHabitForm(w http.ResponseWriter, r *http.Request) error {
	habitID, err := h.getIDFromURL(r)
	if err != nil {
		return err
	}
	return h.render(w, r, formcomponents.UpdateHabit(habitID, forms.CreateHabitForm{}))
}

func (h *HabitHandler) PostUpdateHabit(w http.ResponseWriter, r *http.Request) error {
	habitID, err := h.getIDFromURL(r)
	if err != nil {
		return err
	}

	var form forms.CreateHabitForm
	if err := h.bindFormData(r, &form); err != nil {
		return err
	}

	err = validate.Struct(&form)
	if err != nil {
		form.FieldErrors = h.parseValidationErrors(err)
		return h.render(w, r, formcomponents.UpdateHabit(habitID, form))
	}

	habit, err := h.Queries.UpdateHabit(r.Context(), model.UpdateHabitParams{
		Name: form.Name,
		ID:   habitID,
	})
	if err != nil {
		return h.handleDBError(err)
	}

	http.Redirect(w, r, fmt.Sprintf("/habits/%d", habit.ID), http.StatusSeeOther)
	return nil
}
