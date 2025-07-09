package apperror

import (
	"fmt"
	"net/http"
	"strings"
)

type AppError struct {
	StatusCode int
	Message    string
	Err        error
}

func New(code int, msg string) AppError {
	return AppError{
		StatusCode: code,
		Message:    msg,
	}
}

func NewWrap(code int, msg string, err error) AppError {
	return AppError{
		StatusCode: code,
		Message:    msg,
		Err:        err,
	}
}

func FromMap(code int, m map[string]string) AppError {
	var b strings.Builder
	for k, v := range m {
		fmt.Fprintf(&b, "%s: %s\n", k, v)
	}

	return AppError{
		StatusCode: code,
		Message:    strings.TrimSuffix(b.String(), "\n"),
	}
}

func (e AppError) Error() string {
	msg := fmt.Sprintf("%d: %s", e.StatusCode, e.Message)
	if e.Err != nil {
		return fmt.Sprintf("%s, %s", msg, e.Err.Error())
	}

	return msg
}

func (e AppError) Unwrap() error {
	return e.Err
}

var (
	ErrDuplicateEmail     = New(http.StatusBadRequest, "A user with this email already exists")
	ErrInvalidCredentials = New(http.StatusUnprocessableEntity, "Email or password is not correct")
)
