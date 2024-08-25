insert into habit (name)
values
	('exercise'),
	('study'),
	('take vitamins');

insert into entry (habit_id, entry_date)
values
	(1, '2024-01-01'),
	(1, '2024-01-02'),
	(2, '2024-01-03'),
	(3, '2024-01-04'),
	(3, '2024-01-05');
