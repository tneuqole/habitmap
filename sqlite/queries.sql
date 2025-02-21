-- name: CreateHabit :one
INSERT INTO habits (name, created_at) VALUES (?, unixepoch()) RETURNING *;

-- name: GetHabit :one
SELECT * FROM habits WHERE id = ? LIMIT 1;

-- name: GetHabits :many
SELECT * FROM habits;

-- name: UpdateHabit :one
UPDATE habits SET name = ? WHERE id = ? RETURNING *;

-- name: DeleteHabit :exec
DELETE FROM habits WHERE id = ?;

-- name: CreateEntry :one
INSERT INTO entries (entry_date, habit_id) VALUES (?, ?) RETURNING *;

-- name: GetEntriesForHabit :many
SELECT * FROM entries WHERE habit_id = ? ORDER BY entry_date ASC;

-- name: DeleteEntry :exec
DELETE FROM entries WHERE id = ?;
