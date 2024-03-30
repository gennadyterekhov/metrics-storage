package storage

import (
	"context"
	"github.com/gennadyterekhov/metrics-storage/internal/server/config"
)

type StorageInterface interface {
	Clear()

	AddCounter(key string, value int64)
	SetGauge(key string, value float64)

	GetGauge(name string) (float64, error)
	GetCounter(name string) (int64, error)
	GetGaugeOrZero(name string) float64
	GetCounterOrZero(name string) int64
	GetAllGauges() map[string]float64
	GetAllCounters() map[string]int64

	IsEqual(anotherStorage StorageInterface) (eq bool)

	SaveToDisk(filename string) (err error)
	LoadFromDisk(filename string) (err error)

	IsDB() bool
	GetDB() *DBStorage
	SetContext(ctx context.Context)

	CloseDB() error
}

var MetricsRepository = CreateStorage()

func CreateStorage() StorageInterface {
	if config.Conf.DBDsn == "" {
		return CreateRAMStorage()
	}
	return CreateDBStorage()
}
