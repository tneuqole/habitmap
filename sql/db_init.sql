create table habit (
	id serial primary key,
	name varchar(100)
);

create table entry (
	id serial primary key,
	habit_id int,
	entry_date date,

	foreign key (habit_id) references habit(id) on delete cascade,
	unique(habit_id, entry_date)
);

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
	(3, '2024-01-05'),
	(1, '2024-01-06'),
	(1, '2024-01-07'),
	(1, '2024-01-08'),
	(1, '2024-01-09'),
	(1, '2024-01-10'),
	(1, '2024-01-11'),
	(1, '2024-01-12'),
	(1, '2024-01-13'),
	(1, '2024-01-14'),
	(1, '2024-01-15'),
	(1, '2024-01-20'),
	(1, '2024-01-22'),
	(1, '2024-01-23'),
	(1, '2024-01-24'),
	(1, '2024-01-27'),
	(1, '2024-01-28'),
	(1, '2024-01-30'),
	(1, '2024-01-31');

