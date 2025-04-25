package handlers

import (
	"fmt"
	"net/http"

	"github.com/tneuqole/habitmap/internal/model"
)

type EntryHandler struct {
	*BaseHandler
}

func NewEntryHandler(bh *BaseHandler) *EntryHandler {
	return &EntryHandler{
		BaseHandler: bh,
	}
}

type createEntryForm struct {
	HabitID   int64  `form:"habit_id" validate:"required,notblank"`
	EntryDate string `form:"entry_date" validate:"required,notblank"`
}

func (h *EntryHandler) PostEntry(w http.ResponseWriter, r *http.Request) error {
	var form createEntryForm
	if err := h.bindFormData(r, &form); err != nil {
		return err
	}

	err := validate.Struct(&form)
	if err != nil {
		return err
	}

	params := model.CreateEntryParams{EntryDate: form.EntryDate, HabitID: form.HabitID}
	entry, err := h.Queries.CreateEntry(r.Context(), params)
	if err != nil {
		return err
	}

	http.Redirect(w, r, fmt.Sprintf("/habits/%d", entry.HabitID), http.StatusSeeOther)
	return nil
}

func (h *EntryHandler) DeleteEntry(w http.ResponseWriter, r *http.Request) error {
	entryID, err := h.getIDFromURL(r)
	if err != nil {
		return err
	}

	err = h.Queries.DeleteEntry(r.Context(), entryID)
	if err != nil {
		return err
	}

	w.Header().Set("HX-Redirect", "/habits") // TODO
	w.WriteHeader(http.StatusNoContent)

	return nil
}
