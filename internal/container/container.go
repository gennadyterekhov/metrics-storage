package container

import (
	"github.com/gennadyterekhov/metrics-storage/internal/server/storage"
)

var MetricsRepository = storage.CreateStorage()
