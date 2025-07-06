package ctxutil

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/tneuqole/habitmap/internal/apperror"
)

type contextKey string

const (
	AppErrorCtxKey  contextKey = "appError"
	RequestIDCtxKey contextKey = "requestID"
	LoggerCtxKey    contextKey = "logger"
	NonceCtxKey     contextKey = "nonce"
)

func SetAppError(r *http.Request, err apperror.AppError) *http.Request {
	ctx := r.Context()
	ctx = context.WithValue(ctx, AppErrorCtxKey, err)
	return r.WithContext(ctx)
}

func GetAppError(ctx context.Context) *apperror.AppError {
	if v, ok := ctx.Value(AppErrorCtxKey).(apperror.AppError); ok {
		return &v
	}
	return nil
}

func SetRequestID(r *http.Request, requestID string) *http.Request {
	ctx := r.Context()
	ctx = context.WithValue(ctx, RequestIDCtxKey, requestID)
	return r.WithContext(ctx)
}

func GetRequestID(ctx context.Context) string {
	if v, ok := ctx.Value(RequestIDCtxKey).(string); ok {
		return v
	}
	return ""
}

func SetLogger(r *http.Request, logger *slog.Logger) *http.Request {
	ctx := r.Context()
	ctx = context.WithValue(ctx, LoggerCtxKey, logger)
	return r.WithContext(ctx)
}

func GetLogger(ctx context.Context) *slog.Logger {
	if v, ok := ctx.Value(RequestIDCtxKey).(*slog.Logger); ok {
		return v
	}
	return slog.Default()
}

func SetNonce(r *http.Request, nonce string) *http.Request {
	ctx := r.Context()
	ctx = context.WithValue(ctx, NonceCtxKey, nonce)
	return r.WithContext(ctx)
}

func GetNonce(ctx context.Context) string {
	if v, ok := ctx.Value(NonceCtxKey).(string); ok {
		return v
	}
	return ""
}
