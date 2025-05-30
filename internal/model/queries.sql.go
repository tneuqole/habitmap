// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: queries.sql

package model

import (
	"context"
)

const createEntry = `-- name: CreateEntry :one
INSERT INTO entries (entry_date, habit_id) VALUES (?, ?) RETURNING id, entry_date, habit_id
`

type CreateEntryParams struct {
	EntryDate string
	HabitID   int64
}

func (q *Queries) CreateEntry(ctx context.Context, arg CreateEntryParams) (Entry, error) {
	row := q.db.QueryRowContext(ctx, createEntry, arg.EntryDate, arg.HabitID)
	var i Entry
	err := row.Scan(&i.ID, &i.EntryDate, &i.HabitID)
	return i, err
}

const createHabit = `-- name: CreateHabit :one
INSERT INTO habits (name, created_at) VALUES (?, unixepoch()) RETURNING id, name, created_at
`

func (q *Queries) CreateHabit(ctx context.Context, name string) (Habit, error) {
	row := q.db.QueryRowContext(ctx, createHabit, name)
	var i Habit
	err := row.Scan(&i.ID, &i.Name, &i.CreatedAt)
	return i, err
}

const deleteEntry = `-- name: DeleteEntry :one
DELETE FROM entries
WHERE id = ? RETURNING id, entry_date, habit_id
`

func (q *Queries) DeleteEntry(ctx context.Context, id int64) (Entry, error) {
	row := q.db.QueryRowContext(ctx, deleteEntry, id)
	var i Entry
	err := row.Scan(&i.ID, &i.EntryDate, &i.HabitID)
	return i, err
}

const deleteHabit = `-- name: DeleteHabit :exec
DELETE FROM habits
WHERE id = ?
`

func (q *Queries) DeleteHabit(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteHabit, id)
	return err
}

const getEntriesForHabitByYear = `-- name: GetEntriesForHabitByYear :many
SELECT
    id,
    entry_date,
    habit_id
FROM entries
WHERE habit_id = ? AND strftime('%Y', entry_date) = ?
ORDER BY entry_date ASC
`

type GetEntriesForHabitByYearParams struct {
	HabitID   int64
	EntryDate string
}

func (q *Queries) GetEntriesForHabitByYear(ctx context.Context, arg GetEntriesForHabitByYearParams) ([]Entry, error) {
	rows, err := q.db.QueryContext(ctx, getEntriesForHabitByYear, arg.HabitID, arg.EntryDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Entry
	for rows.Next() {
		var i Entry
		if err := rows.Scan(&i.ID, &i.EntryDate, &i.HabitID); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getEntriesForHabitByYearAndMonth = `-- name: GetEntriesForHabitByYearAndMonth :many
SELECT
    id,
    entry_date,
    habit_id
FROM entries
WHERE habit_id = ? AND strftime('%Y-%m', entry_date) = ?
ORDER BY entry_date ASC
`

type GetEntriesForHabitByYearAndMonthParams struct {
	HabitID   int64
	EntryDate string
}

func (q *Queries) GetEntriesForHabitByYearAndMonth(ctx context.Context, arg GetEntriesForHabitByYearAndMonthParams) ([]Entry, error) {
	rows, err := q.db.QueryContext(ctx, getEntriesForHabitByYearAndMonth, arg.HabitID, arg.EntryDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Entry
	for rows.Next() {
		var i Entry
		if err := rows.Scan(&i.ID, &i.EntryDate, &i.HabitID); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getHabit = `-- name: GetHabit :one
SELECT
    id,
    name,
    created_at
FROM habits
WHERE id = ? LIMIT 1
`

func (q *Queries) GetHabit(ctx context.Context, id int64) (Habit, error) {
	row := q.db.QueryRowContext(ctx, getHabit, id)
	var i Habit
	err := row.Scan(&i.ID, &i.Name, &i.CreatedAt)
	return i, err
}

const getHabits = `-- name: GetHabits :many
SELECT
    id,
    name,
    created_at
FROM habits
`

func (q *Queries) GetHabits(ctx context.Context) ([]Habit, error) {
	rows, err := q.db.QueryContext(ctx, getHabits)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Habit
	for rows.Next() {
		var i Habit
		if err := rows.Scan(&i.ID, &i.Name, &i.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateHabit = `-- name: UpdateHabit :one
UPDATE habits SET name = ?
WHERE id = ?
RETURNING id, name, created_at
`

type UpdateHabitParams struct {
	Name string
	ID   int64
}

func (q *Queries) UpdateHabit(ctx context.Context, arg UpdateHabitParams) (Habit, error) {
	row := q.db.QueryRowContext(ctx, updateHabit, arg.Name, arg.ID)
	var i Habit
	err := row.Scan(&i.ID, &i.Name, &i.CreatedAt)
	return i, err
}
