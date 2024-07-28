package services

import (
	"github.com/gennadyterekhov/metrics-storage/internal/server/repositories"
)

type PingService struct {
	Repository repositories.RepositoryInterface
}

func NewPingService(repo repositories.RepositoryInterface) PingService {
	return PingService{
		Repository: repo,
	}
}
