package database

import (
	"context"

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

func (db Database) GetEntriesByDateRange(ctx context.Context, queryParams model.EntryDateRangeQuery) ([]model.Entry, error) {
	rows, err := db.Conn.Query(context.TODO(), "SELECT * FROM entry WHERE habit_id=$1 AND entry_date BETWEEN $2 AND $3", queryParams.HabitID, queryParams.StartDate, queryParams.EndDate)
	if err != nil {
		return nil, err
	}

	var entries []model.Entry
	entries, err = pgx.CollectRows(rows, pgx.RowToStructByPos[model.Entry])
	return entries, err
}
