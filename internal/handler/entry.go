package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/tneuqole/habitmap/internal/database"
	"github.com/tneuqole/habitmap/internal/model"
)

type EntryHandler struct {
	DB database.Database
}

func (h EntryHandler) PostEntry(c echo.Context) error {
	// TODO validate Content-Type
	var entry model.Entry
	if err := c.Bind(&entry); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	log.Printf("Writing entry to database: %+v", entry)

	// TODO sanitize/validate input
	err := h.DB.CreateEntry(context.TODO(), &entry)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error reading from database: %s", err))
	}

	log.Printf("Wrote entry: %+v\n", entry)
	return c.JSONPretty(http.StatusCreated, entry, "  ")
}

func (h EntryHandler) DeleteEntry(c echo.Context) error {
	id := c.Param("id")
	err := h.DB.DeleteEntry(context.TODO(), id)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error deleting entry: %s", err))
	}

	return c.NoContent(http.StatusNoContent)
}

func (h EntryHandler) GetEntries(c echo.Context) error {
	id, _ := strconv.Atoi(c.QueryParam("habitId"))
	entries, err := h.DB.GetAllEntries(context.TODO(), id)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error collecting entries: %s", err))
	}

	return c.JSONPretty(http.StatusOK, entries, "  ")

}
