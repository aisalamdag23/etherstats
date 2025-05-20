package middleware

import (
	"net/http"

	"github.com/aisalamdag23/etherstats/internal/infrastructure/logger"
	"go.uber.org/zap"
)

// CtxWithLogger is a middleware that puts a logger instance to context
func CtxWithLogger(loggerEntry *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := logger.ToContext(r.Context(), loggerEntry)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
