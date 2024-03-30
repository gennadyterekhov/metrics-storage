package middleware

import (
	"github.com/gennadyterekhov/metrics-storage/internal/server/storage"
	"net/http"
)

func HTTPRequestContextSetter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {

		storage.MetricsRepository.SetContext(request.Context())

		next.ServeHTTP(response, request)
	})
}
