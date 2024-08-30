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
