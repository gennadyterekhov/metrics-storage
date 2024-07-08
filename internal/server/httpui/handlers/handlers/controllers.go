package handlers

import (
	"github.com/gennadyterekhov/metrics-storage/internal/server/services/services"
)

type Controllers struct {
	GetController  GetController
	SaveController SaveController
	PingController PingController
}

func NewControllers(servs *services.Services) Controllers {
	return Controllers{
		GetController:  NewGetController(servs.GetMetricService),
		SaveController: NewSaveController(servs.SaveMetricService),
		PingController: NewPingController(servs.PingService),
	}
}
