package handlers

import (
	"github.com/gennadyterekhov/metrics-storage/internal/server/httpui/middleware"
	"github.com/gennadyterekhov/metrics-storage/internal/server/httpui/middleware/logger"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func GetRouter() chi.Router {

	router := chi.NewRouter()
	registerRoutes(router)

	return router
}

func GetAllMetricsHandler() func(http.ResponseWriter, *http.Request) {
	return middleware.Conveyor(
		http.HandlerFunc(GetMetric),
		logger.RequestAndResponseLoggerMiddleware,
		middleware.MethodGet,
	).ServeHTTP
}

func GetMetricHandler() func(http.ResponseWriter, *http.Request) {
	return middleware.Conveyor(
		http.HandlerFunc(GetMetric),
		logger.RequestAndResponseLoggerMiddleware,
		middleware.MethodGet,
		middleware.URLParametersToGetMetricsArePresent,
	).ServeHTTP
}

func GetMetricJSONHandler() func(http.ResponseWriter, *http.Request) {
	return middleware.Conveyor(
		http.HandlerFunc(GetMetric),
		logger.RequestAndResponseLoggerMiddleware,
		middleware.MethodPost,
	).ServeHTTP
}

func SaveMetricHandler() func(http.ResponseWriter, *http.Request) {
	return middleware.Conveyor(
		http.HandlerFunc(SaveMetric),
		logger.RequestAndResponseLoggerMiddleware,
		middleware.MethodPost,
		middleware.URLParametersToSetMetricsArePresent,
	).ServeHTTP
}

func SaveMetricJSONHandler() func(http.ResponseWriter, *http.Request) {
	return middleware.Conveyor(
		http.HandlerFunc(SaveMetric),
		logger.RequestAndResponseLoggerMiddleware,
		middleware.MethodPost,
	).ServeHTTP
}

func registerRoutes(router chi.Router) {
	router.Head("/", HeadHandler)

	router.Get("/", GetAllMetricsHandler())

	router.Get(
		"/value/{metricType}/{metricName}",
		GetMetricHandler(),
	)
	router.Post(
		"/value/",
		GetMetricJSONHandler(),
	)
	router.Post(
		"/value",
		GetMetricJSONHandler(),
	)

	router.Post(
		"/update/{metricType}/{metricName}/{metricValue}",
		SaveMetricHandler(),
	)
	router.Post(
		"/update/",
		SaveMetricJSONHandler(),
	)
	router.Post(
		"/update",
		SaveMetricJSONHandler(),
	)
}
