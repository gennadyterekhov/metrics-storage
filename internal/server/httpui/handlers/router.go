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

	router.Get(
		"/ping",
		Ping,
	)

	router.Get("/", GetAllMetricsHandler().ServeHTTP)

	router.Get(
		"/value/{metricType}/{metricName}",
		GetMetricHandler().ServeHTTP,
	)
	router.Post(
		"/value/",
		GetMetricHandler().ServeHTTP,
	)
	router.Post(
		"/value",
		GetMetricHandler().ServeHTTP,
	)

	router.Post(
		"/update/{metricType}/{metricName}/{metricValue}",
		SaveMetricHandler().ServeHTTP,
	)
	router.Post(
		"/update/",
		SaveMetricHandler().ServeHTTP,
	)
	router.Post(
		"/update",
		SaveMetricHandler().ServeHTTP,
	)
	router.Post(
		"/updates/",
		SaveMetricListHandler().ServeHTTP,
	)
}
