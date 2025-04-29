# include .env

.PHONY: fmt
fmt:
	golangci-lint run --enable-only goimports --fix ./...

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: test
test:
	go test -race -cover ./... -coverprofile=coverage.out

.PHONY: coverage
coverage: test
	go tool cover -html=coverage.out

.PHONY: build
build:
	esbuild assets/*.js --minify --outdir=public
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

.PHONY: sql/fmt
sql/fmt:
	sqlfluff fix --dialect sqlite sqlite/


