-- Insert sample habits
INSERT INTO habits (name) VALUES
('Exercise'),
('Read'),
('Meditate'),
('Journal'),
('Drink Water'),
('Sleep Early'),
('Practice Coding'),
('Walk Outside');

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
