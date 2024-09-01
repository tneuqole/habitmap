package main

import (
	"database/sql"
	"log"

	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tneuqole/habitmap/internal/database"
	"github.com/tneuqole/habitmap/internal/handler"
)

func main() {
	conn, err := sql.Open("sqlite3", "./habitmap.db")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	db := database.Database{Conn: conn}
	e := echo.New()

	homeHandler := handler.HomeHandler{}
	e.GET("/", homeHandler.GetHome)

	habitHandler := handler.HabitHandler{DB: db}
	e.POST("/habit", habitHandler.PostHabit)
	e.GET("/habit/:id", habitHandler.GetHabit)

	entryHandler := handler.EntryHandler{DB: db}
	e.POST("/entry", entryHandler.PostEntry)
	e.DELETE("/entry/:id", entryHandler.DeleteEntry)
	e.GET("/entry", entryHandler.GetEntries)

	e.Logger.Fatal(e.Start(":8080"))
}
