package ctxutil

import (
	"context"
	"net/http"

	"github.com/tneuqole/habitmap/internal/apperror"
	"github.com/tneuqole/habitmap/internal/session"
)

type contextKey string

var AppErrorCtxKey = contextKey("appError")

func SetAppError(r *http.Request, err apperror.AppError) *http.Request {
	ctx := r.Context()
	ctx = context.WithValue(ctx, AppErrorCtxKey, err)
	return r.WithContext(ctx)
}

func GetAppError(ctx context.Context) (apperror.AppError, bool) {
	val, ok := ctx.Value(AppErrorCtxKey).(apperror.AppError)
	return val, ok
}

var FlashCtxKey = contextKey("flash")

func SetFlash(r *http.Request, msg string) *http.Request {
	ctx := r.Context()
	ctx = context.WithValue(ctx, FlashCtxKey, msg)
	return r.WithContext(ctx)
}

func GetFlash(ctx context.Context) (string, bool) {
	val, ok := ctx.Value(AppErrorCtxKey).(string)
	return val, ok
}

var SessionDataCtxKey = contextKey("sessionData")

func SetSessionData(r *http.Request, sessionData session.SessionData) *http.Request {
	ctx := r.Context()
	ctx = context.WithValue(ctx, SessionDataCtxKey, sessionData)
	return r.WithContext(ctx)
}

func GetSessionData(ctx context.Context) (session.SessionData, bool) {
	val, ok := ctx.Value(SessionDataCtxKey).(session.SessionData)
	return val, ok
}
