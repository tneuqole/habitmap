run:
	@npx tailwindcss -i ./assets/app.css -o ./public/app.css
	@templ generate
	@go run main.go

css:
	@npx tailwindcss -i ./assets/app.css -o ./public/app.css --watch
