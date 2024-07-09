package middleware

import (
	"net/http"

	"github.com/gennadyterekhov/metrics-storage/internal/server/config"

	"github.com/gennadyterekhov/metrics-storage/internal/server/httpui/middleware/compressor"
	"github.com/gennadyterekhov/metrics-storage/internal/server/httpui/middleware/hashchecker"
	"github.com/gennadyterekhov/metrics-storage/internal/server/httpui/middleware/logger"
)

type Set struct {
	Config *config.ServerConfig
}

type Middleware func(http.Handler) http.Handler

func New(conf *config.ServerConfig) *Set {
	return &Set{
		Config: conf,
	}
}

func (set *Set) CommonConveyor(h http.Handler, middlewares ...Middleware) http.Handler {
	allMiddlewares := getCommonMiddlewares(set.Config)
	allMiddlewares = append(allMiddlewares, middlewares...)
	return conveyor(h, allMiddlewares...)
}

func getCommonMiddlewares(conf *config.ServerConfig) []Middleware {
	return []Middleware{
		logger.RequestAndResponseLoggerMiddleware,
		compressor.GzipCompressor,
		ContentType,
		hashchecker.New(conf.PayloadSignatureKey).CheckHash,
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
