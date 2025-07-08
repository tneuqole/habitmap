package logutil

import "log/slog"

const (
	ErrorLogKey     = "error"
	RequestIDLogKey = "requestID"
)

func ErrorSlog(err error) slog.Attr {
	return slog.Any(ErrorLogKey, err.Error())
}
