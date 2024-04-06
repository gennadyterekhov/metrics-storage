package hashchecker

import (
	"github.com/gennadyterekhov/metrics-storage/internal/helper/hasher"
	"github.com/gennadyterekhov/metrics-storage/internal/server/config"
	"net/http"
)

func CheckHash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		ok, hash := isHashValid(request)
		if ok {
			response.Header().Set("HashSHA256", hash)
			next.ServeHTTP(response, request)
			return
		}

		if hash == "" {
			next.ServeHTTP(response, request)
			return
		}

		response.WriteHeader(http.StatusBadRequest)
	})
}

func isHashValid(request *http.Request) (bool, string) {
	return hasher.IsBodyHashValid(request, config.Conf.PayloadSignatureKey), request.Header.Get("HashSHA256")
}
