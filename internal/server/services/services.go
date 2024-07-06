package services

import (
	"github.com/gennadyterekhov/metrics-storage/internal/server/config"
	"github.com/gennadyterekhov/metrics-storage/internal/server/repositories"
)

type Services struct {
	GetMetricService  GetMetricService
	SaveMetricService SaveMetricService
	TimeTracker       TimeTracker
}

func NewServices(repo repositories.RepositoryInterface, conf *config.ServerConfig) Services {
	return Services{
		GetMetricService:  NewGetMetricService(repo),
		SaveMetricService: NewSaveMetricService(repo, conf),
		TimeTracker:       NewTimeTracker(repo, conf),
	}
}
