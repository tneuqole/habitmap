package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/tneuqole/heatmap/internal/date"
)

type Entry struct {
	ID      int             `json:"id"`
	HabitID int             `json:"habitId"`
	Date    date.CustomDate `json:"date" binding:"required"`
}

type EntryQuery struct {
	HabitID   int    `query:"habitId"`
	StartDate string `query:"startDate"`
	EndDate   string `query:"endDate"`
}

type EntryHandler struct {
	Conn *pgx.Conn
}

func (h EntryHandler) PostEntry(c echo.Context) error {
	// TODO validate Content-Type
	entry := new(Entry)
	if err := c.Bind(entry); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	log.Printf("Writing entry to database: %+v", entry)

	// TODO sanitize/validate input
	var id int
	err := h.Conn.QueryRow(context.TODO(), "INSERT INTO entry(habit_id, entry_date) VALUES($1, $2) RETURNING id", entry.HabitID, entry.Date).Scan(&id)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error reading from database: %s", err))
	}

	entry.ID = id
	log.Printf("Wrote entry: %+v\n", entry)
	return c.JSONPretty(http.StatusCreated, entry, "  ")
}

func (h EntryHandler) GetEntries(c echo.Context) error {
	var entries []Entry
	var entryQuery EntryQuery

	if err := c.Bind(&entryQuery); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	log.Printf("Got query parameters: %+v\n", entryQuery)

	rows, err := h.Conn.Query(context.TODO(), "SELECT * FROM entry WHERE habit_id=$1 AND entry_date BETWEEN $2 AND $3", entryQuery.HabitID, entryQuery.StartDate, entryQuery.EndDate)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error reading from database: %s", err))
	}

	entries, err = pgx.CollectRows(rows, pgx.RowToStructByPos[Entry])
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error collecting entries: %s", err))
	}

	return c.JSONPretty(http.StatusOK, entries, "  ")

}
