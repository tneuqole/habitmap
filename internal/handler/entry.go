package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"

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

func (h EntryHandler) GetEntries(c echo.Context) error {
	var queryParams model.EntryDateRangeQuery
	if err := c.Bind(&queryParams); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	log.Printf("Got query parameters: %+v\n", queryParams)

	entries, err := h.DB.GetEntriesByDateRange(context.TODO(), queryParams)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error collecting entries: %s", err))
	}

	return c.JSONPretty(http.StatusOK, entries, "  ")

}
