// package main sets up initial configuration and starts the web server
package main

import (
	"database/sql"
	"log/slog"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tneuqole/habitmap/internal/handlers"
	"github.com/tneuqole/habitmap/internal/model"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	db, err := sql.Open("sqlite3", "./habitmap.db") // TODO: probably shouldn't expose filename
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("Error closing database connection: %v", err)
		}
	}()

	queries := model.New(db)

	e := echo.New()
	e.Logger.SetLevel(log.DEBUG) // TODO: make env var
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	e.GET("/health", handlers.GetHealth)

	e.Static("/public", "public")

	baseHandler := &handlers.BaseHandler{
		Logger:  logger,
		Queries: queries,
	}

	habitHandler := handlers.NewHabitHandler(baseHandler)
	e.GET("/habits", habitHandler.GetHabits)
	e.GET("/habits/:id", habitHandler.GetHabit)
	e.DELETE("/habits/:id", habitHandler.DeleteHabit)
	e.GET("/habits/new", habitHandler.GetCreateHabitForm)
	e.POST("/habits/new", habitHandler.PostHabit)
	e.GET("/habits/:id/edit", habitHandler.GetUpdateHabitForm)
	e.POST("/habits/:id/edit", habitHandler.PostUpdateHabit)

	entryHandler := handlers.NewEntryHandler(baseHandler)
	e.POST("/entries", entryHandler.PostEntry)
	e.DELETE("/entries/:id", entryHandler.DeleteEntry)

	e.Logger.Fatal(e.Start(":4000"))
}
