package app

import (
	"github.com/gennadyterekhov/metrics-storage/internal/server/config"
	"github.com/gennadyterekhov/metrics-storage/internal/server/httpui/handlers"
	"github.com/gennadyterekhov/metrics-storage/internal/server/repositories"
	"github.com/gennadyterekhov/metrics-storage/internal/server/services"
)

type App struct {
	Config      config.ServerConfig
	Services    services.Services
	Controllers handlers.Controllers
	Repository  repositories.MetricsRepository
}

func New() App {
	return App{}
}
