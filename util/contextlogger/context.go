package contextlogger

import (
	"context"
	"log/slog"
)

type key int

var loggerKey key

func New(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func From(ctx context.Context) *slog.Logger {
	if ctxLogger := ctx.Value(loggerKey); ctxLogger != nil {
		if logger, ok := ctxLogger.(*slog.Logger); ok {
			return logger
		}
	}

	return slog.New(discardHandler{})
}
