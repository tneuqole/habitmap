CREATE TABLE habits (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    created_at INTEGER NOT NULL
);

CREATE TABLE entries (
    id INTEGER PRIMARY KEY,
    entry_date TEXT NOT NULL CHECK (entry_date LIKE '____-__-__'),
    habit_id INTEGER NOT NULL,
    FOREIGN KEY (habit_id) REFERENCES habits (id) ON DELETE CASCADE
)
