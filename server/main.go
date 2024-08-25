package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

var db *pgx.Conn

type Habit struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Entry struct {
	ID      int    `json:"id"`
	HabitID int    `json:"habitId"`
	Date    string `json:"date"`
}

func postHabit(c *gin.Context) {
	var habit Habit
	if err := c.BindJSON(&habit); err != nil {
		return // TODO
	}

	// TODO: write to db
}

func getHabitByID(c *gin.Context) {
	id := c.Param("id")

	var habit Habit
	err := db.QueryRow(context.TODO(), "SELECT * FROM habit WHERE id = $1", id).Scan(&habit.ID, &habit.Name)
	if err != nil {
		log.Fatal(err)
	}
	c.IndentedJSON(http.StatusOK, habit)
}

func main() {
	ctx := context.Background()
	var err error
	db, err = pgx.Connect(ctx, "postgresql://myuser:password@localhost:5433/heatmap?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close(ctx)
	router := gin.Default()
	router.POST("/habit", postHabit)
	router.GET("/habit/:id", getHabitByID)

	router.Run("localhost:8080")
}
