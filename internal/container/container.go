package container

import "github.com/gennadyterekhov/metrics-storage/internal/storage"

type Container struct {
	MemStorage *storage.MemStorage
}

var Instance = &Container{
	MemStorage: storage.CreateStorage(),
}
