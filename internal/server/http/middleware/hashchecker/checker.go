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
		hash, ok, err := mdl.isHashValid(request)
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
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

func (mdl Hashchecker) isHashValid(request *http.Request) (string, bool, error) {
	hash := request.Header.Get("HashSHA256")
	ok, err := hasher.IsBodyHashValid(request, mdl.PayloadSignatureKey)
	return hash, ok, err
}
