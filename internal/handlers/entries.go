package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tneuqole/habitmap/internal/model"
)

// EntryHandler handles HTTP requests related to entries
type EntryHandler struct {
	*BaseHandler
}

// NewEntryHandler creates a new EntryHandler
func NewEntryHandler(bh *BaseHandler) *EntryHandler {
	return &EntryHandler{
		BaseHandler: bh,
	}
}

type createEntryForm struct {
	HabitID   int64  `form:"habit_id" validate:"required,notblank"`
	EntryDate string `form:"entry_date" validate:"required,notblank"`
}

// PostEntry processes a form for creating a new entry
func (h *EntryHandler) PostEntry(c echo.Context) error {
	form := createEntryForm{}
	if err := c.Bind(&form); err != nil {
		return err
	}

	err := validate.Struct(&form)
	if err != nil {
		return err
	}

	params := model.CreateEntryParams{EntryDate: form.EntryDate, HabitID: form.HabitID}
	entry, err := h.Queries.CreateEntry(c.Request().Context(), params)
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/habits/%d", entry.HabitID))
}

type deleteEntryParams struct {
	EntryID int64 `param:"id"`
}

// DeleteEntry deletes an entry by id
func (h *EntryHandler) DeleteEntry(c echo.Context) error {
	params := deleteEntryParams{}
	if err := c.Bind(&params); err != nil {
		return err
	}

	err := h.Queries.DeleteEntry(c.Request().Context(), params.EntryID)
	if err != nil {
		return err
	}

	c.Response().Header().Add("Hx-Redirect", "/habits") // TODO

	return nil
}
