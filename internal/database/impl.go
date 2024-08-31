package database

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/tneuqole/habitmap/internal/model"
)

type Database struct {
	Conn *pgx.Conn
}

func (db Database) CreateHabit(ctx context.Context, habit *model.Habit) error {
	var id int
	err := db.Conn.QueryRow(context.TODO(), "INSERT INTO habit(name) VALUES($1) RETURNING id", habit.Name).Scan(&id)

	habit.ID = id
	return err
}

func (db Database) GetHabitByID(ctx context.Context, id string) (model.Habit, error) {
	var habit model.Habit
	err := db.Conn.QueryRow(context.TODO(), "SELECT * FROM habit WHERE id = $1", id).Scan(&habit.ID, &habit.Name)

	return habit, err

}

func (db Database) CreateEntry(ctx context.Context, entry *model.Entry) error {
	var id int
	err := db.Conn.QueryRow(ctx, "INSERT INTO entry(habit_id, entry_date) VALUES($1, $2) RETURNING id", entry.HabitID, entry.Date).Scan(&id)

	entry.ID = id
	return err
}

func (db Database) DeleteEntry(ctx context.Context, id string) error {
	_, err := db.Conn.Exec(context.TODO(), "DELETE FROM entry WHERE id = $1", id)
	return err
}

func (db Database) GetEntriesByDateRange(ctx context.Context, habitID int, startDate, endDate time.Time) ([]model.Entry, error) {
	rows, err := db.Conn.Query(context.TODO(), "SELECT * FROM entry WHERE habit_id=$1 AND entry_date BETWEEN $2 AND $3", habitID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	return collectEntries(rows)
}

func (db Database) GetAllEntries(ctx context.Context, habitID int) ([]model.Entry, error) {
	rows, err := db.Conn.Query(context.TODO(), "SELECT * FROM entry WHERE habit_id=$1", habitID)
	if err != nil {
		return nil, err
	}

	return collectEntries(rows)
}

func collectEntries(rows pgx.Rows) ([]model.Entry, error) {
	return pgx.CollectRows(rows, pgx.RowToStructByPos[model.Entry])
}
