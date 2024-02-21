package handlers

import (
	"github.com/gennadyterekhov/metrics-storage/internal/middleware"
	"github.com/gennadyterekhov/metrics-storage/internal/middleware/logger"
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
		middleware.MethodGet,
		logger.RequestAndResponseLoggerMiddleware,
	).ServeHTTP
}

func GetMetricHandler() func(http.ResponseWriter, *http.Request) {
	return middleware.Conveyor(
		http.HandlerFunc(GetMetric),
		middleware.MethodGet,
		middleware.URLParametersToGetMetricsArePresent,
		logger.RequestAndResponseLoggerMiddleware,
	).ServeHTTP
}

func SaveMetricHandler() func(http.ResponseWriter, *http.Request) {
	return middleware.Conveyor(
		http.HandlerFunc(SaveMetric),
		middleware.MethodPost,
		middleware.URLParametersToSetMetricsArePresent,
		logger.RequestAndResponseLoggerMiddleware,
	).ServeHTTP
}

func registerRoutes(router chi.Router) {
	router.Get("/", GetAllMetricsHandler())
	router.Get(
		"/value/{metricType}/{metricName}",
		GetMetricHandler(),
	)
	router.Post(
		"/update/{metricType}/{metricName}/{metricValue}",
		SaveMetricHandler(),
	)
}
