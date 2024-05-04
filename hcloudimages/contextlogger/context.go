package contextlogger

import (
	"context"
	"log/slog"
)

type key int

var loggerKey key

// New saves the logger as a value to the context. This can then be retrieved through [From].
func New(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// From returns the [*slog.Logger] set on the context by [New]. If there is none,
// it returns a no-op logger that discards any output it receives.
func From(ctx context.Context) *slog.Logger {
	if ctxLogger := ctx.Value(loggerKey); ctxLogger != nil {
		if logger, ok := ctxLogger.(*slog.Logger); ok {
			return logger
		}
	}

	return slog.New(discardHandler{})
}
