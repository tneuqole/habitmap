package handlers

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/tneuqole/habitmap/internal/model"
	"github.com/tneuqole/habitmap/internal/templates"
	"github.com/tneuqole/habitmap/internal/templates/forms"
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

func (h *HabitHandler) GetHabit(w http.ResponseWriter, r *http.Request) error {
	habitID, err := h.getIDFromURL(r)
	if err != nil {
		return err
	}

	habit, err := h.Queries.GetHabit(r.Context(), habitID)
	if err != nil {
		return h.handleDBError(err)
	}

	entries, err := h.Queries.GetEntriesForHabitByYear(r.Context(), model.GetEntriesForHabitByYearParams{
		HabitID:   habitID,
		EntryDate: "2025",
	})
	if err != nil {
		return h.handleDBError(err)
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
	return h.render(w, r, pages.Habit(habit, sortedMonths, entriesForMonths))
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
	return h.render(w, r, forms.CreateHabit(templates.HabitFormData{}))
}

type createHabitForm struct {
	Name string `form:"name" validate:"required,notblank,min=1,max=32"`
}

func (h *HabitHandler) PostHabit(w http.ResponseWriter, r *http.Request) error {
	var form createHabitForm
	if err := h.bindFormData(r, &form); err != nil {
		return err
	}

	err := validate.Struct(&form)
	if err != nil {
		errors := parseValidationErrors(err)
		data := templates.HabitFormData{
			Name:   form.Name,
			Errors: errors,
		}
		return h.render(w, r, forms.CreateHabit(data))
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
	return h.render(w, r, forms.UpdateHabit(templates.HabitFormData{ID: habitID}))
}

func (h *HabitHandler) PostUpdateHabit(w http.ResponseWriter, r *http.Request) error {
	habitID, err := h.getIDFromURL(r)
	if err != nil {
		return err
	}

	var form createHabitForm
	if err := h.bindFormData(r, &form); err != nil {
		return err
	}

	err = validate.Struct(&form)
	if err != nil {
		errors := parseValidationErrors(err)
		data := templates.HabitFormData{
			Name:   form.Name,
			Errors: errors,
		}
		return h.render(w, r, forms.UpdateHabit(data))
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
