package logium

import (
	"context"
	"net/http"
)

type Logger interface {
	With(args ...any) Logger
	WithField(key string, value any) Logger
	WithFields(fields map[string]any) Logger
	WithError(err error) Logger
	WithRequest(r *http.Request) Logger
	WithOperation(operation string) Logger
	WithComponent(component string) Logger

	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)

	DebugContext(ctx context.Context, msg string, args ...any)
	InfoContext(ctx context.Context, msg string, args ...any)
	WarnContext(ctx context.Context, msg string, args ...any)
	ErrorContext(ctx context.Context, msg string, args ...any)
}
