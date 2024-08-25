create table habit (
	id serial primary key,
	name varchar(100)
);

create table entry (
	id serial primary key,
	habit_id int,
	entry_date date,

	foreign key (habit_id) references habit(id) on delete cascade
)
