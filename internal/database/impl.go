package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/tneuqole/habitmap/internal/model"
)

type Database struct {
	Conn *sql.DB
}

func (db Database) CreateHabit(ctx context.Context, habit *model.Habit) error {
	res, err := db.Conn.Exec("INSERT INTO habit(name) VALUES($1)", habit.Name)
	if err != nil {
		return err
	}

	id, _ := res.LastInsertId()
	habit.ID = int(id)

	return nil
}

func (db Database) GetHabitByID(ctx context.Context, id string) (model.Habit, error) {
	var habit model.Habit
	err := db.Conn.QueryRow("SELECT * FROM habit WHERE id = $1", id).Scan(&habit.ID, &habit.Name)

	return habit, err

}

func (db Database) CreateEntry(ctx context.Context, entry *model.Entry) error {
	res, err := db.Conn.Exec("INSERT INTO entry(habit_id, entry_date) VALUES($1, $2) RETURNING id", entry.HabitID, entry.Date)
	if err != nil {
		return err
	}

	id, _ := res.LastInsertId()
	entry.ID = int(id)

	return nil
}

func (db Database) DeleteEntry(ctx context.Context, id string) (model.Entry, error) {
	entry := model.NewEntry()
	err := db.Conn.QueryRow("DELETE FROM entry WHERE id = $1 RETURNING habit_id, entry_date", id).Scan(&entry.HabitID, &entry.Date)
	return entry, err
}

func (db Database) GetEntriesByDateRange(ctx context.Context, habitID int, startDate, endDate time.Time) ([]model.Entry, error) {
	rows, err := db.Conn.Query("SELECT * FROM entry WHERE habit_id=$1 AND entry_date BETWEEN $2 AND $3", habitID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	return collectEntries(rows)
}

func (db Database) GetAllEntries(ctx context.Context, habitID int) ([]model.Entry, error) {
	rows, err := db.Conn.Query("SELECT * FROM entry WHERE habit_id=$1", habitID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	return collectEntries(rows)
}

func collectEntries(rows *sql.Rows) ([]model.Entry, error) {
	var entries []model.Entry
	for rows.Next() {
		var entry model.Entry
		err := rows.Scan(&entry.ID, &entry.HabitID, &entry.Date)
		if err != nil {
			return nil, err
		}

		entries = append(entries, entry)
	}

	return entries, nil
}
