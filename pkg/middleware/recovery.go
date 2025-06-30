package middleware

import (
	"net/http"

	"go.uber.org/zap"
)

func RecoverPanic(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Error("Recovered from panic",
						zap.Any("error", err),
						zap.String("method", r.Method),
						zap.String("path", r.URL.Path),
					)
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"error":"internal server error"}`))
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
