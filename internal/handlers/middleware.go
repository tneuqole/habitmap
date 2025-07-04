package handlers

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/tneuqole/habitmap/internal/apperror"
	"github.com/tneuqole/habitmap/internal/util"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

func (h *BaseHandler) Wrap(f HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			h.Logger.Error("API_ERROR", util.ErrorSlog(err), slog.String("error_type", fmt.Sprintf("%T", err)))

			var appErr apperror.AppError
			if errors.As(err, &appErr) {
				w.WriteHeader(appErr.StatusCode)
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

func (h *BaseHandler) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !h.Session.Data(r.Context()).IsAuthenticated {
			http.Redirect(w, r, "/users/login", http.StatusSeeOther)
			return
		}

		// set cache-control to no-store so pages that require auth are not cached
		w.Header().Add("cache-control", "no-store")
		next.ServeHTTP(w, r)
	})
}
