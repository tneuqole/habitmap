package util

import (
	"fmt"
	"strings"
)

type AppError struct {
	StatusCode int
	Message    string
}

func NewAppError(code int, msg string) *AppError {
	return &AppError{
		StatusCode: code,
		Message:    msg,
	}
}

func NewAppErrorFromMap(code int, m map[string]string) *AppError {
	var b strings.Builder
	for k, v := range m {
		fmt.Fprintf(&b, "%s: %s\n", k, v)
	}

	return &AppError{
		StatusCode: code,
		Message:    b.String(),
	}
}

func (e AppError) Error() string {
	return fmt.Sprintf("%d: %s", e.StatusCode, e.Message)
}
