package tools

import (
	"context"
	"log/slog"
)

var ctxKeyLogger struct{} = struct{}{}

func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, ctxKeyLogger, logger)
}

func LoggerFromContext(ctx context.Context) *slog.Logger {
	if value, ok := ctx.Value(ctxKeyLogger).(*slog.Logger); !ok {
		if value == nil {
			panic("context is nil")
		} else {
			panic("context is not *slog.Logger, also is not nil")
		}
	} else {
		return value
	}
}
