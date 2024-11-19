package main

import (
	"database/sql"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tneuqole/habitmap/internal/handlers"
)

func main() {
	conn, err := sql.Open("sqlite3", "./habitmap.db")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	e := echo.New()
	e.Use(middleware.Logger())

	e.Static("/public", "public")

	e.GET("/health", handlers.GetHealth)

	e.Logger.Fatal(e.Start(":8080"))
}
