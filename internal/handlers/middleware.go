package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/tneuqole/habitmap/internal/apperror"
	"github.com/tneuqole/habitmap/internal/ctxutil"
	"github.com/tneuqole/habitmap/internal/logutil"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

func (h *BaseHandler) Wrap(f HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := ctxutil.GetLogger(r.Context())

		if err := f(w, r); err != nil {
			logger.Error("API_ERROR", logutil.ErrorSlog(err), slog.String("errorType", fmt.Sprintf("%T", err)))

			statusCode := http.StatusInternalServerError
			var appErr apperror.AppError
			if errors.As(err, &appErr) {
				statusCode = appErr.StatusCode
			}

			h.RenderErrorPage(w, r, statusCode)
		}
	}
}

func (h *BaseHandler) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !h.Session.Data(r.Context()).IsAuthenticated {
			http.Redirect(w, r, "/users/login", http.StatusSeeOther)
			return
		}

		// don't cache pages that require auth
		w.Header().Add("Expires", "Thu, 01 Jan 1970 00:00:00 UTC")
		w.Header().Add("Cache-Control", "no-cache, private, max-age=0")
		w.Header().Add("X-Accel-Expires", "0")
		w.Header().Add("Pragma", "no-cache")

		next.ServeHTTP(w, r)
	})
}

func (h *BaseHandler) SetHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf := make([]byte, 16) //nolint:mnd
		_, _ = rand.Read(buf)   //nolint:errcheck
		nonce := base64.StdEncoding.EncodeToString(buf)

		r = ctxutil.SetNonce(r, nonce)

		// set a strict CSP to prevent XSS, clickjacking, code injection.
		w.Header().Set("Content-Security-Policy", fmt.Sprintf(
			"default-src 'self'; "+
				"script-src 'self' https://unpkg.com 'nonce-%s'; "+
				"style-src 'self' fonts.googleapis.com 'nonce-%s'; "+
				"font-src fonts.gstatic.com;",
			nonce, nonce,
		))

		// include full URL for same site requests.
		// for all other requests path & query params are removed.
		w.Header().Add("Referrer-Policy", "origin-when-cross-origin")

		// tell browsers to not sniff MIME type.
		// helps prevent content sniffing attacks.
		w.Header().Add("X-Content-Type-Options", "nosniff")

		// prevent clickjacking in older browsers.
		w.Header().Add("X-Frame-Options", "deny")

		// recommended to disable because of CSP header.
		w.Header().Add("X-XSS-Protection", "0")

		next.ServeHTTP(w, r)
	})
}

type customResponseWriter struct {
	http.ResponseWriter
	statusCode  int
	bytes       int
	wroteHeader bool
}

func (w *customResponseWriter) WriteHeader(statusCode int) {
	if !w.wroteHeader {
		w.statusCode = statusCode
		w.ResponseWriter.WriteHeader(statusCode)
		w.wroteHeader = true
	}
}

func (w *customResponseWriter) Write(buf []byte) (int, error) {
	// default to 200 OK if statusCode was not written
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	n, err := w.ResponseWriter.Write(buf)
	w.bytes += n
	return n, err
}

// LogRequest generates a unique request id, attaches it to a logger,
// and stores the logger in ctx for usage while processing the request.
// It also logs information about the request and response.
func (h *BaseHandler) LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		requestID := uuid.NewString()
		r = ctxutil.SetRequestID(r, requestID)

		logger := slog.Default().With(slog.String(logutil.RequestIDLogKey, requestID))
		r = ctxutil.SetLogger(r, logger)

		ww := &customResponseWriter{ResponseWriter: w}

		defer func() {
			duration := time.Since(start)
			logger.Info(
				"REQUEST_INFO",
				slog.String("ip", r.RemoteAddr),
				slog.String("proto", r.Proto),
				slog.String("method", r.Method),
				slog.String("uri", r.RequestURI),
				slog.Int("statusCode", ww.statusCode),
				slog.Duration("duration", duration),
				slog.Int("bytes", ww.bytes),
			)
		}()

		next.ServeHTTP(ww, r)
	})
}

func (h *BaseHandler) RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				middleware.PrintPrettyStack(err)
				w.Header().Set("Connection", "close")
				h.RenderErrorPage(w, r, http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// Timeout is a simple middleware that sets a timeout on the
// request context. It is only helpful for operations that respect
// context deadlines like database queries.
func (h *BaseHandler) Timeout(timeout time.Duration) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
