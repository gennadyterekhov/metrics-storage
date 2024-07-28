package storage

import (
	"context"
)

type Interface interface {
	Clear()

	AddCounter(ctx context.Context, key string, value int64)
	SetGauge(ctx context.Context, key string, value float64)

	GetGauge(ctx context.Context, name string) (float64, error)
	GetCounter(ctx context.Context, name string) (int64, error)
	GetGaugeOrZero(ctx context.Context, name string) float64
	GetCounterOrZero(ctx context.Context, name string) int64
	GetAllGauges(ctx context.Context) map[string]float64
	GetAllCounters(ctx context.Context) map[string]int64

	SaveToDisk(ctx context.Context, filename string) (err error)
	LoadFromDisk(ctx context.Context, filename string) (err error)

	GetDB() *DBStorage
	GetMemStorage() *MemStorage

	CloseDB() error
}

func New(dsn string) Interface {
	if dsn == "" {
		return NewRAMStorage()
	}
	return NewDBStorage(dsn)
}
