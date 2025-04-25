# include .env

.PHONY: build
build:
	npx tailwindcss -i ./assets/app.css -o ./public/app.css
	templ generate

.PHONY: dev
dev:
	air

.PHONY: run
run:
	go run ./cmd/api

.PHONY: sql
sql:
	sqlite3 habitmap.db

.PHONY: sql/init
sql/init:
	sqlite3 habitmap.db < sqlite/schema.sql
	sqlite3 habitmap.db < sqlite/seed.sql

.PHONY: sql/gen
sql/gen:
	sqlc generate

