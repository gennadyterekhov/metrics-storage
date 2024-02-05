package container

import (
	"github.com/gennadyterekhov/metrics-storage/internal/repositories"
	"github.com/gennadyterekhov/metrics-storage/internal/storage"
)

type Container struct {
	MetricsRepository repositories.MetricsRepository
}

var Instance = &Container{
	MetricsRepository: storage.CreateStorage(),
}
