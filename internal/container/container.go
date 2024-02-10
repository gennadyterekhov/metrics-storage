package container

import (
	"github.com/gennadyterekhov/metrics-storage/internal/storage"
)

var MetricsRepository = storage.CreateStorage()
