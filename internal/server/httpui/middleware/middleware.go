package middleware

import (
	"github.com/gennadyterekhov/metrics-storage/internal/server/httpui/middleware/compressor"
	"github.com/gennadyterekhov/metrics-storage/internal/server/httpui/middleware/logger"
	"net/http"
)

type Middleware func(http.Handler) http.Handler

func CommonConveyor(h http.Handler, middlewares ...Middleware) http.Handler {
	allMiddlewares := getCommonMiddlewares()
	allMiddlewares = append(allMiddlewares, middlewares...)
	return conveyor(h, allMiddlewares...)
}

func getCommonMiddlewares() []Middleware {
	return []Middleware{
		logger.RequestAndResponseLoggerMiddleware,
		compressor.GzipCompressor,
		ContentType,
	}
}

func conveyor(h http.Handler, middlewares ...Middleware) http.Handler {
	middlewaresLength := len(middlewares)
	// in reverse, so that middlewares are applied in order that they are passed in router
	for i := middlewaresLength - 1; i >= 0; i -= 1 {
		h = middlewares[i](h)
	}
	return h
}
