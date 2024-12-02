package handlers

import (
	"fmt"

	"github.com/a-h/templ"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
	"github.com/labstack/echo/v4"
)

func Render(c echo.Context, component templ.Component) error {
	return component.Render(c.Request().Context(), c.Response())
}

func NewValidate() *validator.Validate {
	validate := validator.New()
	validate.RegisterValidation("notblank", validators.NotBlank)
	return validate
}

func ParseValidationErrors(err error) map[string]string {
	errors := make(map[string]string)

	if _, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range err.(validator.ValidationErrors) {
			var msg string
			switch fieldErr.Tag() {
			case "required":
				msg = fmt.Sprintf("%s is required", fieldErr.Field())
			case "min":
				msg = fmt.Sprintf("%s must be at least %s characters long", fieldErr.Field(), fieldErr.Param())
			case "max":
				msg = fmt.Sprintf("%s must be at most %s characters long", fieldErr.Field(), fieldErr.Param())
			case "alpha":
			case "alphanum":
			case "alphanumunicode":
			case "alphaunicode":
			case "ascii":
				msg = fmt.Sprintf("%s contains invalid characters", fieldErr.Field())
			case "notblank":
				msg = fmt.Sprintf("%s cannot be blank", fieldErr.Field())
			default:
				msg = fmt.Sprintf("%s is invalid", fieldErr.Field())
			}
			errors[fieldErr.Field()] = msg
		}
	}

	return errors
}
