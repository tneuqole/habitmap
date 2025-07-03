INSERT INTO users (name, email, hashed_password)
VALUES ('user', 'user@habitmap.com', '$2a$10$HfCZNLQaTHEz3zGofuodA.CHx2yw5jXFOgAU5bfo8R9AYcz/qOcDW');

INSERT INTO habits (user_id, name) VALUES
(1, 'Exercise'),
(1, 'Read'),
(1, 'Meditate'),
(1, 'Journal'),
(1, 'Drink Water'),
(1, 'Sleep Early'),
(1, 'Practice Coding'),
(1, 'Walk Outside');

-- Generate habit entries for 2 years (~10 per month per habit)
WITH RECURSIVE date_series AS (
    SELECT date('now', '-2 years', 'start of month') AS entry_date
    UNION ALL
    SELECT date(entry_date, '+1 day') FROM date_series
    WHERE entry_date < date('now', 'start of month', '-1 day')
)

INSERT INTO entries (entry_date, habit_id)
SELECT
    ds.entry_date,
    h.id
FROM date_series AS ds
INNER JOIN habits AS h
    ON (strftime('%d', ds.entry_date) % 3 = 0) -- Roughly every 3rd day
ORDER BY ds.entry_date, h.id;
