build:
	@npx tailwindcss -i ./assets/app.css -o ./public/app.css
	@templ generate

