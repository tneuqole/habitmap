CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    hashed_password TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE habits (
    id INTEGER PRIMARY KEY,
    -- user_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    created_at INTEGER NOT NULL
    -- created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP

    -- FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE entries (
    id INTEGER PRIMARY KEY,
    entry_date TEXT NOT NULL CHECK (entry_date LIKE '____-__-__'),
    habit_id INTEGER NOT NULL,

    FOREIGN KEY (habit_id) REFERENCES habits (id) ON DELETE CASCADE
);

-- ALTER TABLE entries ADD COLUMN year TEXT GENERATED ALWAYS AS (strftime('%Y', entry_date)) STORED;
-- ALTER TABLE entries ADD COLUMN year_month TEXT GENERATED ALWAYS AS (strftime('%Y-%m', entry_date)) STORED;
--
-- CREATE INDEX idx_entries_year ON entries(year);
-- CREATE INDEX idx_entries_year_month ON entries(year_month);

