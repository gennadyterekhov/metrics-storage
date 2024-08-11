package middleware

import (
	"net/http"

	"github.com/gennadyterekhov/metrics-storage/internal/common/constants"
	"github.com/gennadyterekhov/metrics-storage/internal/common/logger"
)

func ContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request != nil && request.Header.Get(constants.HeaderContentType) == constants.ApplicationJSON {
			logger.Custom.Debugln("ContentType set to json ")
			response.Header().Set(constants.HeaderContentType, constants.ApplicationJSON)
		} else {
			logger.Custom.Debugln("ContentType set to text/html")
			response.Header().Set(constants.HeaderContentType, constants.TextHTML)
		}

		next.ServeHTTP(response, request)
	})
}
