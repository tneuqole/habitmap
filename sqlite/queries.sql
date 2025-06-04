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
INSERT INTO habits (name, created_at)
VALUES (?, unixepoch())
RETURNING *;

-- name: GetHabit :one
SELECT
    id,
    name,
    created_at
FROM habits
WHERE id = ? LIMIT 1;

-- name: GetHabits :many
SELECT
    id,
    name,
    created_at
FROM habits;

-- name: UpdateHabit :one
UPDATE habits SET name = ?
WHERE id = ?
RETURNING *;

-- name: DeleteHabit :exec
DELETE FROM habits
WHERE id = ?;

-- name: CreateEntry :one
INSERT INTO entries
(entry_date, habit_id)
VALUES (?, ?)
RETURNING *;

-- name: GetEntriesForHabitByYear :many
SELECT
    id,
    entry_date,
    habit_id
FROM entries
WHERE habit_id = ? AND strftime('%Y', entry_date) = ?
ORDER BY entry_date ASC;

-- name: GetEntriesForHabitByYearAndMonth :many
SELECT
    id,
    entry_date,
    habit_id
FROM entries
WHERE habit_id = ? AND strftime('%Y-%m', entry_date) = ?
ORDER BY entry_date ASC;

-- name: DeleteEntry :one
DELETE FROM entries
WHERE id = ? RETURNING *;
