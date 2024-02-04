package router

import (
	"github.com/gennadyterekhov/metrics-storage/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func GetRouter() chi.Router {
	router := chi.NewRouter()
	registerRoutes(router)

	return router
}

func registerRoutes(router chi.Router) {
	router.Get("/", handlers.GetAllMetrics)
	router.Get("/value/{metricType}/{metricName}", handlers.GetMetric)
	router.Post("/update/{metricType}/{metricName}/{metricValue}", handlers.SaveMetric)
}
