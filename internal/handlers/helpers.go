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

func NewValidate() (*validator.Validate, error) {
	validate := validator.New()
	err := validate.RegisterValidation("notblank", validators.NotBlank)
	if err != nil {
		return nil, err
	}

	return validate, nil
}

func ParseValidationErrors(err any) map[string]string {
	formErrors := make(map[string]string)

	if valErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range valErrors {
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
			formErrors[fieldErr.Field()] = msg
		}
	}

	return formErrors
}
