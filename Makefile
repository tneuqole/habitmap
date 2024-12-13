include .env

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

.PHONY: lint
lint:
	golangci-lint run --fix ./...


.PHONY: test
test:
	go test -race -cover ./... -coverprofile=coverage.out

.PHONY: cov
cov: test
	go tool cover -html=coverage.out

.PHONY: sql
sql:
	sqlite3 habitmap.db

.PHONY: sql/init
sql/init:
	sqlite3 habitmap.db < sqlite/schema.sql

.PHONY: sql/gen
sql/gen:
	sqlc generate

