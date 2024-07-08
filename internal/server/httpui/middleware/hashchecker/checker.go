package hashchecker

import (
	"net/http"

	"github.com/gennadyterekhov/metrics-storage/internal/common/helper/hasher"
)

type Hashchecker struct {
	PayloadSignatureKey string
}

func New(payloadSignatureKey string) Hashchecker {
	return Hashchecker{
		PayloadSignatureKey: payloadSignatureKey,
	}
}

func (mdl Hashchecker) CheckHash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		ok, hash := mdl.isHashValid(request)
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

func (mdl Hashchecker) isHashValid(request *http.Request) (bool, string) {
	return hasher.IsBodyHashValid(request, mdl.PayloadSignatureKey), request.Header.Get("HashSHA256")
}
