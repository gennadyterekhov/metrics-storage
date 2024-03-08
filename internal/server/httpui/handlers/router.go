package handlers

import (
	"github.com/go-chi/chi/v5"
)

func GetRouter() chi.Router {
	router := chi.NewRouter()
	registerRoutes(router)

	return router
}

func registerRoutes(router chi.Router) {
	router.Head("/", HeadHandler)

	router.Get("/", GetAllMetricsHandlerFunc())

	router.Get(
		"/value/{metricType}/{metricName}",
		GetMetricHandlerFunc(),
	)
	router.Post(
		"/value/",
		GetMetricHandlerFunc(),
	)
	router.Post(
		"/value",
		GetMetricHandlerFunc(),
	)

	router.Post(
		"/update/{metricType}/{metricName}/{metricValue}",
		SaveMetricHandlerFunc(),
	)
	router.Post(
		"/update/",
		SaveMetricHandlerFunc(),
	)
	router.Post(
		"/update",
		SaveMetricHandlerFunc(),
	)
	router.Post(
		"/update/batch",
		SaveMetricBatchHandlerFunc(),
	)
	router.Post(
		"/update/batch/",
		SaveMetricBatchHandlerFunc(),
	)
}
