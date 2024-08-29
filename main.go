package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/tneuqole/heatmap/handler"
)

func main() {
	ctx := context.TODO()
	conn, err := pgx.Connect(ctx, "postgresql://myuser:password@localhost:5433/heatmap?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(ctx)
	e := echo.New()

	homeHandler := handler.HomeHandler{}
	e.GET("/", homeHandler.GetHome)

	habitHandler := handler.HabitHandler{Conn: conn}
	e.POST("/habit", habitHandler.PostHabit)
	e.GET("/habit/:id", habitHandler.GetHabitByID)

	entryHandler := handler.EntryHandler{Conn: conn}
	e.POST("/entry", entryHandler.PostEntry)
	e.GET("/entry", entryHandler.GetEntries)

	e.Logger.Fatal(e.Start(":8080"))
}
