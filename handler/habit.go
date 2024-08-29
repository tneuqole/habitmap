package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

type Habit struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type HabitHandler struct {
	Conn *pgx.Conn
}

func (h HabitHandler) PostHabit(c echo.Context) error {
	// TODO validate Content-Type
	var habit Habit
	if err := c.Bind(&habit); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	// TODO sanitize/validate input
	var id int
	err := h.Conn.QueryRow(context.TODO(), "INSERT INTO habit(name) VALUES($1) RETURNING id", habit.Name).Scan(&id)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error reading from database: %s", err))
	}

	habit.ID = id
	log.Printf("Wrote habit: %+v\n", habit)
	return c.JSONPretty(http.StatusCreated, habit, "  ")
}

func (h HabitHandler) GetHabitByID(c echo.Context) error {
	id := c.Param("id")

	var habit Habit
	err := h.Conn.QueryRow(context.TODO(), "SELECT * FROM habit WHERE id = $1", id).Scan(&habit.ID, &habit.Name)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error reading from database: %s", err))
	}

	return c.JSONPretty(http.StatusOK, habit, "  ")
}
