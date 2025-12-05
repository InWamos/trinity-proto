package middleware

import (
	"log/slog"
	"net/http"
)

type LoggingMiddleware struct {
	logger *slog.Logger
}

func NewLoggingMiddleware(logger *slog.Logger) *LoggingMiddleware {
	middlewareLogger := logger.With(slog.String("component", "logging_middleware"))
	return &LoggingMiddleware{logger: middlewareLogger}
}

func (middleware *LoggingMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		middleware.logger.InfoContext(r.Context(),
			"request",
			slog.String("http_method", r.Method),
			slog.String("http_uri", r.RequestURI),
			slog.String("remore_addr", r.RemoteAddr))
		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
