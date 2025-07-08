package handlers

import (
	"net/http"

	"github.com/tneuqole/habitmap/internal/forms"
	"github.com/tneuqole/habitmap/internal/model"
	"github.com/tneuqole/habitmap/internal/templates/components"
)

type EntryHandler struct {
	*BaseHandler
}

func NewEntryHandler(bh *BaseHandler) *EntryHandler {
	return &EntryHandler{
		BaseHandler: bh,
	}
}

func (h *EntryHandler) PostEntry(w http.ResponseWriter, r *http.Request) error {
	var form forms.CreateEntryForm
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
		return h.handleDBError(r.Context(), err)
	}

	return h.render(w, r, components.Entry(entry))
}

func (h *EntryHandler) DeleteEntry(w http.ResponseWriter, r *http.Request) error {
	entryID, err := h.getIDFromURL(r)
	if err != nil {
		return err
	}

	entry, err := h.Queries.DeleteEntry(r.Context(), entryID)
	if err != nil {
		return h.handleDBError(r.Context(), err)
	}
	entry.ID = 0

	return h.render(w, r, components.Entry(entry))
}
