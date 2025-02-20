package main

import (
	"database/sql"
	"flag"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tneuqole/habitmap/internal/handlers"
	"github.com/tneuqole/habitmap/internal/model"
)

func main() {
	dbFile := flag.String("db", "habitmap.db", "sqlite database file")
	flag.Parse()

	validate, err := handlers.NewValidate()
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("sqlite3", *dbFile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	queries := model.New(db)

	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	e.GET("/health", handlers.GetHealth)

	e.Static("/public", "public")

	habitsHandler := handlers.NewHabitsHandler(queries, validate)
	e.GET("/habits", habitsHandler.GetHabits)
	e.GET("/habits/:id", habitsHandler.GetHabit)
	e.DELETE("/habits/:id", habitsHandler.DeleteHabit)
	e.GET("/habits/new", habitsHandler.GetCreateHabitForm)
	e.POST("/habits/new", habitsHandler.PostHabit)
	e.GET("/habits/:id/edit", habitsHandler.GetUpdateHabitForm)
	e.POST("/habits/:id/edit", habitsHandler.PostUpdateHabit)

	e.Logger.Fatal(e.Start(":4000"))
}
