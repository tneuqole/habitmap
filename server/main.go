package main

import (
	"context"
	"database/sql/driver"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

var db *pgx.Conn

const DateFormat = "2006-01-02"

type CustomDate time.Time

// used by gin for dates in POST body
func (d *CustomDate) UnmarshalJSON(data []byte) error {
	// handle empty values?
	parsed, err := time.Parse(`"`+DateFormat+`"`, string(data))
	if err != nil {
		return err
	}

	*d = CustomDate(parsed)
	return nil
}

// used by gin for response body
func (d CustomDate) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(DateFormat)+2)
	b = append(b, '"')
	b = time.Time(d).AppendFormat(b, DateFormat)
	b = append(b, '"')
	return b, nil
}

// used by pgx when writing to db
func (d CustomDate) Value() (driver.Value, error) {
	if d.String() == "0001-01-01 00:00:00" {
		return nil, nil
	}

	return []byte(time.Time(d).Format(DateFormat)), nil
}

// used by pgx when reading from db
func (d *CustomDate) Scan(v interface{}) error {
	// parse from postgres date format
	parsed, _ := time.Parse("2006-01-02 00:00:00 +0000 UTC", v.(time.Time).String())
	*d = CustomDate(parsed)
	return nil
}

func (d CustomDate) String() string {
	return time.Time(d).Format(DateFormat)
}

type Habit struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Entry struct {
	ID      int        `json:"id"`
	HabitID int        `json:"habitId"`
	Date    CustomDate `json:"date" binding:"required"`
}

type EntryQuery struct {
	HabitID   int    `form:"habitId"`
	StartDate string `form:"startDate"`
	EndDate   string `form:"endDate"`
}

func postHabit(c *gin.Context) {
	var habit Habit
	if err := c.BindJSON(&habit); err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// TODO sanitize/validate input
	var id int
	err := db.QueryRow(context.TODO(), "INSERT INTO habit(name) VALUES($1) RETURNING id", habit.Name).Scan(&id)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	habit.ID = id
	log.Printf("Wrote habit: %+v\n", habit)
	c.IndentedJSON(http.StatusCreated, habit)
}

func getHabitByID(c *gin.Context) {
	id := c.Param("id")

	var habit Habit
	err := db.QueryRow(context.TODO(), "SELECT * FROM habit WHERE id = $1", id).Scan(&habit.ID, &habit.Name)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.IndentedJSON(http.StatusOK, habit)
}

func postEntry(c *gin.Context) {
	var entry Entry
	if err := c.BindJSON(&entry); err != nil {
		log.Printf("Error binding json: %s\n", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// TODO sanitize/validate input
	var id int
	err := db.QueryRow(context.TODO(), "INSERT INTO entry(habit_id, entry_date) VALUES($1, $2) RETURNING id", entry.HabitID, entry.Date).Scan(&id)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	entry.ID = id
	log.Printf("Wrote entry: %+v\n", entry)
	c.IndentedJSON(http.StatusCreated, entry)
}

func getEntries(c *gin.Context) {
	var entries []Entry
	var entryQuery EntryQuery

	if err := c.BindQuery(&entryQuery); err != nil {
		log.Printf("Error binding query parameters: %s\n", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	log.Printf("Got query parameters: %+v\n", entryQuery)

	rows, err := db.Query(context.TODO(), "SELECT * FROM entry WHERE habit_id=$1 AND entry_date BETWEEN $2 AND $3", entryQuery.HabitID, entryQuery.StartDate, entryQuery.EndDate)
	if err != nil {
		log.Printf("Error querying entries: %s\n", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	entries, err = pgx.CollectRows(rows, pgx.RowToStructByPos[Entry])
	if err != nil {
		log.Printf("Error collecting entries: %s\n", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.IndentedJSON(http.StatusOK, entries)

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

	router.POST("/entry", postEntry)
	router.GET("/entry", getEntries)

	router.Run("localhost:8080")
}
