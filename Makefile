include .env

build:
	npx tailwindcss -i ./assets/app.css -o ./public/app.css
	templ generate

dev:
	air

run:
	go run ./cmd/api

sql:
	sqlite3 habitmap.db

sql/init:
	sqlite3 habitmap.db < sqlite/schema.sql

