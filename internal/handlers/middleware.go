package handlers

import (
	"errors"
	"log/slog"
	"net/http"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

func (h *BaseHandler) Wrap(f HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			h.Logger.Error("API_ERROR", slog.Any("error", err.Error()))

			var appErr AppError
			if errors.As(err, &appErr) {
				w.WriteHeader(http.StatusOK)
				// TODO: render error page
				w.Write([]byte(appErr.Error())) //nolint:errcheck,gosec
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				// TODO: render error page
				w.Write([]byte("Internal Server Error")) //nolint:errcheck,gosec
			}
		}
	}
}
