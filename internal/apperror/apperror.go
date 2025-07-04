package apperror

import (
	"fmt"
	"net/http"
	"strings"
)

type AppError struct {
	StatusCode int
	Message    string
}

func New(code int, msg string) AppError {
	return AppError{
		StatusCode: code,
		Message:    msg,
	}
}

func FromMap(code int, m map[string]string) AppError {
	var b strings.Builder
	for k, v := range m {
		fmt.Fprintf(&b, "%s: %s\n", k, v)
	}

	return AppError{
		StatusCode: code,
		Message:    b.String(),
	}
}

func (e AppError) Error() string {
	return fmt.Sprintf("%d: %s", e.StatusCode, e.Message)
}

var (
	ErrDuplicateEmail     = New(http.StatusBadRequest, "A user with this email already exists")
	ErrInvalidCredentials = New(http.StatusUnprocessableEntity, "Email or password is not correct")
)
