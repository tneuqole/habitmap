package ctxutil

import (
	"context"
	"net/http"

	"github.com/tneuqole/habitmap/internal/apperror"
)

type contextKey string

var AppErrorCtxKey = contextKey("appError")

func SetAppError(r *http.Request, err apperror.AppError) *http.Request {
	ctx := r.Context()
	ctx = context.WithValue(ctx, AppErrorCtxKey, err)
	return r.WithContext(ctx)
}

func GetAppError(ctx context.Context) (apperror.AppError, bool) {
	val := ctx.Value(AppErrorCtxKey)
	appErr, ok := val.(apperror.AppError)
	return appErr, ok
}
