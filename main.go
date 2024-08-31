package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/tneuqole/habitmap/internal/database"
	"github.com/tneuqole/habitmap/internal/handler"
)

func main() {
	ctx := context.TODO()
	conn, err := pgx.Connect(ctx, "postgresql://myuser:password@localhost:5433/habitmap?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(ctx)

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
