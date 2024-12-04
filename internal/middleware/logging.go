package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

// LoggingMiddleware создает middleware для логирования HTTP запросов
func LoggingMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Создаем ResponseWriter, который может отслеживать статус ответа
			wrappedWriter := &responseWriter{
				ResponseWriter: w,
				status:         http.StatusOK,
			}

			// Логируем входящий запрос
			logger.Info("Incoming request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("remote_addr", r.RemoteAddr),
				zap.String("user_agent", r.UserAgent()),
			)

			// Передаем запрос дальше
			next.ServeHTTP(wrappedWriter, r)

			// Логируем результат запроса
			duration := time.Since(start)
			logger.Info("Request completed",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", wrappedWriter.status),
				zap.Duration("duration", duration),
			)
		})
	}
}

// responseWriter оборачивает http.ResponseWriter для отслеживания статуса ответа
type responseWriter struct {
	http.ResponseWriter
	status int
}

// WriteHeader обрабатывает статус ответа
func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
