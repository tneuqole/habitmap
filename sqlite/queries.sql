-- name: CreateUser :one
INSERT INTO users (name, email, hashed_password)
VALUES (?, ?, ?)
RETURNING
    id;

-- name: GetUser :one
SELECT
    id,
    name,
    email,
    hashed_password
FROM users
WHERE email = ?;

-- name: GetUserByID :one
SELECT
    name,
    email,
    created_at
FROM users
WHERE id = ?;

-- name: CreateHabit :one
INSERT INTO habits (name)
VALUES (?)
RETURNING id;

-- name: GetHabit :one
SELECT * FROM habits
WHERE id = ? LIMIT 1;

-- name: GetHabits :many
SELECT * FROM habits;

-- name: UpdateHabit :exec
UPDATE habits
SET
    name = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: DeleteHabit :exec
DELETE FROM habits
WHERE id = ?;

-- name: CreateEntry :one
INSERT INTO entries
(entry_date, habit_id)
VALUES (?, ?)
RETURNING *;

-- name: GetEntriesForHabitByYear :many
SELECT * FROM entries
WHERE habit_id = ? AND year = ?
ORDER BY entry_date ASC;

-- name: GetEntriesForHabitByYearAndMonth :many
SELECT * FROM entries
WHERE habit_id = ? AND year_month = ?
ORDER BY entry_date ASC;

-- name: DeleteEntry :one
DELETE FROM entries
WHERE id = ?
RETURNING *;
