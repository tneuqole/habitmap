package main

import (
	"database/sql"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tneuqole/habitmap/internal/handlers"
	"github.com/tneuqole/habitmap/internal/model"
)

func main() {
	db, err := sql.Open("sqlite3", "./habitmap.db") // TODO: probably shouldn't expose filename
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

	habitHandler := handlers.NewHabitHandler(queries)
	e.GET("/habits", habitHandler.GetHabits)
	e.GET("/habits/:id", habitHandler.GetHabit)
	e.DELETE("/habits/:id", habitHandler.DeleteHabit)
	e.GET("/habits/new", habitHandler.GetCreateHabitForm)
	e.POST("/habits/new", habitHandler.PostHabit)
	e.GET("/habits/:id/edit", habitHandler.GetUpdateHabitForm)
	e.POST("/habits/:id/edit", habitHandler.PostUpdateHabit)

	entryHandler := handlers.NewEntryHandler(queries)
	e.POST("/entries", entryHandler.PostEntry)

	e.Logger.Fatal(e.Start(":4000"))
}
