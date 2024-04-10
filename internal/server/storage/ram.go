package storage

import (
	"context"
	"fmt"
	"github.com/gennadyterekhov/metrics-storage/internal/constants/exceptions"
)

type MemStorage struct {
	Counters           map[string]int64   `json:"counters"`
	Gauges             map[string]float64 `json:"gauges"`
	HTTPRequestContext context.Context    `json:"-"`
}

func CreateRAMStorage() *MemStorage {
	return &MemStorage{
		Counters: make(map[string]int64, 0),
		Gauges:   make(map[string]float64, 0),
	}
}

func (strg *MemStorage) Clear() {
	strg.Counters = make(map[string]int64, 0)
	strg.Gauges = make(map[string]float64, 0)
}

func (strg *MemStorage) hasGauge(name string) bool {
	_, ok := strg.Gauges[name]
	return ok
}

func (strg *MemStorage) hasCounter(name string) bool {
	_, ok := strg.Counters[name]
	return ok
}

func (strg *MemStorage) AddCounter(ctx context.Context, key string, value int64) {
	strg.Counters[key] += value
}

func (strg *MemStorage) SetGauge(ctx context.Context, key string, value float64) {
	strg.Gauges[key] = value
}

func (strg *MemStorage) GetGauge(ctx context.Context, name string) (float64, error) {
	if !strg.hasGauge(name) {
		return 0, fmt.Errorf(exceptions.UnknownMetricName)
	}
	return strg.GetGaugeOrZero(ctx, name), nil
}

func (strg *MemStorage) GetCounter(ctx context.Context, name string) (int64, error) {
	if !strg.hasCounter(name) {
		return 0, fmt.Errorf(exceptions.UnknownMetricName)
	}
	return strg.GetCounterOrZero(ctx, name), nil
}

func (strg *MemStorage) GetGaugeOrZero(ctx context.Context, name string) float64 {
	val, ok := strg.Gauges[name]
	if !ok {
		return 0
	}
	return val
}

func (strg *MemStorage) GetCounterOrZero(ctx context.Context, name string) int64 {
	val, ok := strg.Counters[name]
	if !ok {
		return 0
	}
	return val
}

func (strg *MemStorage) GetAllGauges(ctx context.Context) map[string]float64 {
	return strg.Gauges
}

func (strg *MemStorage) GetAllCounters(ctx context.Context) map[string]int64 {
	return strg.Counters
}

func (strg *MemStorage) CloseDB() error {
	return nil
}

func (strg *MemStorage) GetDB() *DBStorage {
	return nil
}

func (strg *MemStorage) GetMemStorage() *MemStorage {
	return strg
}
