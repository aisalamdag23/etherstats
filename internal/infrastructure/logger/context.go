package logger

import (
	"context"

	"go.uber.org/zap"
)

type ctxLoggerMarker struct{}

type ctxLogger struct {
	logger *zap.Logger
}

var (
	ctxLoggerKey = &ctxLoggerMarker{}
)

func newCtxLogger(entry *zap.Logger) *ctxLogger {
	return &ctxLogger{
		logger: entry,
	}
}

// ToContext adds the logger to the context for extraction later.
// Returning the new context that has been created.
func ToContext(ctx context.Context, entry *zap.Logger) context.Context {
	return context.WithValue(ctx, ctxLoggerKey, newCtxLogger(entry))
}

// Extract takes the call-scoped logger from context.
func Extract(ctx context.Context) *zap.Logger {
	if ctx == nil {
		return newLogger(zap.DebugLevel)
	}
	l := extract(ctx)
	if l == nil {
		return newLogger(zap.DebugLevel)
	}
	return l.logger
}

func extract(ctx context.Context) *ctxLogger {
	if ctx == nil {
		return nil
	}
	l, ok := ctx.Value(ctxLoggerKey).(*ctxLogger)
	if !ok || l == nil || l.logger == nil {
		return nil
	}
	return l
}
