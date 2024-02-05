package handlers

import (
	"github.com/gennadyterekhov/metrics-storage/internal/middleware"
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
	).ServeHTTP
}

func GetMetricHandler() func(http.ResponseWriter, *http.Request) {
	return middleware.Conveyor(
		http.HandlerFunc(GetMetric),
		middleware.MethodGet,
		middleware.URLParametersToGetMetricsArePresent,
	).ServeHTTP
}

func SaveMetricHandler() func(http.ResponseWriter, *http.Request) {
	return middleware.Conveyor(
		http.HandlerFunc(SaveMetric),
		middleware.MethodPost,
		middleware.URLParametersToSetMetricsArePresent,
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
