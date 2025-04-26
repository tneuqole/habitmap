package util

import "log/slog"

func ErrorSlog(err error) slog.Attr {
	return slog.Any("error", err.Error())
}
