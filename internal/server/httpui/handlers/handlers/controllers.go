package handlers

import (
	"github.com/gennadyterekhov/metrics-storage/internal/server/httpui/middleware"
	"github.com/gennadyterekhov/metrics-storage/internal/server/services/services"
)

//	@title			metrics-storage API
//	@version		1.0
//	@description	metrics-storage API.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	https://github.com/gennadyterekhov/metrics-storage/issues/
//	@contact.email	mail@example.com

//	@license.name	MIT
//	@license.url	https://en.wikipedia.org/wiki/MIT_License

//	@host		localhost:8080
//	@BasePath	/

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
