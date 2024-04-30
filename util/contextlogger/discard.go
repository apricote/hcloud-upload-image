package contextlogger

import (
	"context"
	"log/slog"
)

// discardHandler is a [slog.Handler] that just discards any input. It is a safe default if any library
// method does not get passed a logger through the context.
type discardHandler struct{}

func (discardHandler) Enabled(_ context.Context, _ slog.Level) bool  { return false }
func (discardHandler) Handle(_ context.Context, _ slog.Record) error { return nil }
func (d discardHandler) WithAttrs(_ []slog.Attr) slog.Handler        { return d }
func (d discardHandler) WithGroup(_ string) slog.Handler             { return d }
