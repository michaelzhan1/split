package logs

import (
	"log/slog"
	"net/http"
	"time"
)

func RequestLogger(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            // Use a ResponseWriter wrapper to capture status code
            rw := &responseWriter{w, http.StatusOK}
            next.ServeHTTP(rw, r)
            duration := time.Since(start)

            logger.Info("HTTP request",
                slog.String("method", r.Method),
                slog.String("path", r.URL.Path),
                slog.Int("status", rw.statusCode),
                slog.Duration("duration", duration),
                slog.String("remote_addr", r.RemoteAddr),
            )
        })
	}
}

type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}
