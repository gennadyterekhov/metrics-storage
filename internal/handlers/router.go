package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
)

func GetRouter() chi.Router {
	fmt.Println("func GetRouter")

	router := chi.NewRouter()
	registerRoutes(router)

	return router
}

func registerRoutes(router chi.Router) {
	router.Get("/", GetAllMetrics)
	router.Get("/value/{metricType}/{metricName}", GetMetric)
	router.Post("/update/{metricType}/{metricName}/{metricValue}", SaveMetric)
}
