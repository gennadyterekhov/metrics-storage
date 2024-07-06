package router

import (
	"github.com/gennadyterekhov/metrics-storage/internal/server/httpui/handlers"
	"github.com/go-chi/chi/v5"
)

type Router struct {
	Controllers *handlers.Controllers
	ChiRouter   chi.Router
}

func New(conts *handlers.Controllers) Router {
	chiRouter := chi.NewRouter()

	instance := Router{
		Controllers: conts,
		ChiRouter:   chiRouter,
	}
	instance.registerRoutes()

	return instance
}

// deprecated
func GetRouter() chi.Router {
	router := chi.NewRouter()
	// registerRoutes(router)

	return router
}

func (rtr Router) registerRoutes() {
	rtr.ChiRouter.Head("/", handlers.HeadHandler)

	rtr.ChiRouter.Get(
		"/ping",
		handlers.Ping,
	)

	rtr.ChiRouter.Get("/", handlers.GetAllMetricsHandler(rtr.Controllers.GetController).ServeHTTP)

	rtr.ChiRouter.Get(
		"/value/{metricType}/{metricName}",
		handlers.GetMetricHandler(rtr.Controllers.GetController).ServeHTTP,
	)
	rtr.ChiRouter.Post(
		"/value/",
		handlers.GetMetricHandler(rtr.Controllers.GetController).ServeHTTP,
	)
	rtr.ChiRouter.Post(
		"/value",
		handlers.GetMetricHandler(rtr.Controllers.GetController).ServeHTTP,
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
		handlers.SaveMetricHandler(rtr.Controllers.SaveController).ServeHTTP,
	)
	rtr.ChiRouter.Post(
		"/update",
		handlers.SaveMetricHandler(rtr.Controllers.SaveController).ServeHTTP,
	)
}
