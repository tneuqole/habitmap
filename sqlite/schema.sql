CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    hashed_password TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE habits (
    id INTEGER PRIMARY KEY,
    user_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    UNIQUE (user_id, name)
);

CREATE TABLE entries (
    id INTEGER PRIMARY KEY,
    entry_date TEXT NOT NULL CHECK (entry_date GLOB '[0-9][0-9][0-9][0-9]-[0-1][0-9]-[0-3][0-9]'),
    habit_id INTEGER NOT NULL,

    FOREIGN KEY (habit_id) REFERENCES habits (id) ON DELETE CASCADE,
    UNIQUE (habit_id, entry_date)
);

ALTER TABLE entries ADD COLUMN year TEXT GENERATED ALWAYS AS (strftime('%Y', entry_date)) STORED;
ALTER TABLE entries ADD COLUMN year_month TEXT GENERATED ALWAYS AS (strftime('%Y-%m', entry_date)) STORED;

CREATE INDEX idx_entries_year ON entries(year);
CREATE INDEX idx_entries_year_month ON entries(year_month);

