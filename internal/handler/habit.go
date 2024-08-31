package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/tneuqole/habitmap/internal/database"
	"github.com/tneuqole/habitmap/internal/model"
	"github.com/tneuqole/habitmap/internal/template"
	"github.com/tneuqole/habitmap/internal/util"
)

type HabitHandler struct {
	DB database.Database
}

func (h HabitHandler) PostHabit(c echo.Context) error {
	// TODO validate Content-Type
	var habit model.Habit
	if err := c.Bind(&habit); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	// TODO sanitize/validate input
	err := h.DB.CreateHabit(context.TODO(), &habit)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error reading from database: %s", err))
	}

	log.Printf("Wrote habit: %+v\n", habit)
	return c.JSONPretty(http.StatusCreated, habit, "  ")
}

func (h HabitHandler) GetHabitByID(c echo.Context) error {
	id := c.Param("id")

	habit, err := h.DB.GetHabitByID(context.TODO(), id)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error reading from database: %s", err))
	}

	return c.JSONPretty(http.StatusOK, habit, "  ")
}

type QueryParams struct {
	HabitID int `param:"id"`
	Year    int `query:"year"`
	Month   int `query:"month"`
}

func (h HabitHandler) GetHabit(c echo.Context) error {
	var params QueryParams
	if err := c.Bind(&params); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	if params.Month != 0 && params.Year == 0 {
		return c.String(http.StatusBadRequest, "If specifying month, must specify year")
	}

	var err error
	var entries []model.Entry
	if params.Year == 0 {
		entries, err = h.DB.GetAllEntries(context.TODO(), params.HabitID)
	} else if params.Month != 0 {
		startDate := time.Date(params.Year, time.Month(params.Month), 1, 0, 0, 0, 0, time.UTC)
		endDate := startDate.AddDate(0, 1, -1)
		entries, err = h.DB.GetEntriesByDateRange(context.TODO(), params.HabitID, startDate, endDate)
	} else {
		startDate := time.Date(params.Year, time.January, 1, 0, 0, 0, 0, time.UTC)
		endDate := startDate.AddDate(1, 0, -1)
		entries, err = h.DB.GetEntriesByDateRange(context.TODO(), params.HabitID, startDate, endDate)
	}

	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error reading from database: %s", err))
	}

	return util.Render(c, template.Month(entries[0].Date.Time(), entries))
}
