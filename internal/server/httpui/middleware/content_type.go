package middleware

import (
	"github.com/gennadyterekhov/metrics-storage/internal/constants"
	"github.com/gennadyterekhov/metrics-storage/internal/logger"
	"net/http"
)

func ContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request != nil && request.Header.Get(constants.HeaderContentType) == constants.ApplicationJSON {
			logger.ZapSugarLogger.Debugln("ContentType set to json ")
			response.Header().Set(constants.HeaderContentType, constants.ApplicationJSON)
		} else {
			logger.ZapSugarLogger.Debugln("ContentType set to text/html")
			response.Header().Set(constants.HeaderContentType, constants.TextHTML)
		}

		next.ServeHTTP(response, request)
	})
}
