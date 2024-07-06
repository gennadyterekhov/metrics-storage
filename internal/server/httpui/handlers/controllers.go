package handlers

import (
	"github.com/gennadyterekhov/metrics-storage/internal/server/services"
)

type Controllers struct {
	GetController  GetController
	SaveController SaveController
}

func NewControllers(servs *services.Services) Controllers {
	return Controllers{
		GetController:  NewGetController(servs.GetMetricService),
		SaveController: NewSaveController(servs.SaveMetricService),
	}
}
