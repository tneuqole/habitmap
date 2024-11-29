CREATE TABLE habits (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    created_at INTEGER NOT NULL
);

CREATE TABLE entries (
    id INTEGER PRIMARY KEY,
    entry_date INTEGER NOT NULL,
    habit_id INTEGER NOT NULL,
    FOREIGN KEY (habit_id) REFERENCES habits (id) ON DELETE CASCADE
)
