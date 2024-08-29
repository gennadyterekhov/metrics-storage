package router

import (
	"github.com/gennadyterekhov/metrics-storage/internal/server/http/handlers/handlers"
	"github.com/go-chi/chi/v5"
)

type Router struct {
	Controllers *handlers.Controllers
	ChiRouter   chi.Router
}

func New(conts *handlers.Controllers) *Router {
	chiRouter := chi.NewRouter()

	instance := Router{
		Controllers: conts,
		ChiRouter:   chiRouter,
	}
	instance.registerRoutes()

	return &instance
}

func (rtr Router) registerRoutes() {
	rtr.ChiRouter.Get(
		"/ping",
		handlers.PingHandler(rtr.Controllers.PingController).ServeHTTP,
	)

	rtr.ChiRouter.Get("/", handlers.GetAllMetricsHandler(rtr.Controllers.GetController).ServeHTTP)

	rtr.ChiRouter.Get(
		"/value/{metricType}/{metricName}",
		handlers.GetMetricHandler(rtr.Controllers.GetController).ServeHTTP,
	)
	rtr.ChiRouter.Post(
		"/value/",
		handlers.GetMetricJSONHandler(rtr.Controllers.GetController).ServeHTTP,
	)
	rtr.ChiRouter.Post(
		"/value",
		handlers.GetMetricJSONHandler(rtr.Controllers.GetController).ServeHTTP,
	)

	rtr.ChiRouter.Post(
		"/updates",
		handlers.SaveMetricListHandler(rtr.Controllers.SaveController).ServeHTTP,
	)
	rtr.ChiRouter.Post(
		"/updates/",
		handlers.SaveMetricListHandler(rtr.Controllers.SaveController).ServeHTTP,
	)

	rtr.ChiRouter.Post(
		"/update/{metricType}/{metricName}/{metricValue}",
		handlers.SaveMetricHandler(rtr.Controllers.SaveController).ServeHTTP,
	)
	rtr.ChiRouter.Post(
		"/update/",
		handlers.SaveMetricJSONHandler(rtr.Controllers.SaveController).ServeHTTP,
	)
	rtr.ChiRouter.Post(
		"/update",
		handlers.SaveMetricJSONHandler(rtr.Controllers.SaveController).ServeHTTP,
	)
}
