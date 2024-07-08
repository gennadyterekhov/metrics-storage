package handlers

import (
	"github.com/gennadyterekhov/metrics-storage/internal/server/httpui/middleware"
	"github.com/gennadyterekhov/metrics-storage/internal/server/services/services"
)

type Controllers struct {
	GetController  GetController
	SaveController SaveController
	PingController PingController
	MiddlewareSet  *middleware.Set
}

func NewControllers(servs *services.Services, middlewareSet *middleware.Set) Controllers {
	return Controllers{
		GetController:  NewGetController(servs.GetMetricService, middlewareSet),
		SaveController: NewSaveController(servs.SaveMetricService, middlewareSet),
		PingController: NewPingController(servs.PingService),
		MiddlewareSet:  middlewareSet,
	}
}
