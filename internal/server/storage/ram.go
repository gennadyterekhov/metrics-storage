package storage

import (
	"context"
	"fmt"
	"sync"

	"github.com/gennadyterekhov/metrics-storage/internal/common/constants/exceptions"
)

type MemStorage struct {
	Counters           map[string]int64   `json:"counters"`
	Gauges             map[string]float64 `json:"gauges"`
	HTTPRequestContext context.Context    `json:"-"`
	mu                 sync.Mutex
}

func NewRAMStorage() *MemStorage {
	return &MemStorage{
		Counters: make(map[string]int64, 0),
		Gauges:   make(map[string]float64, 0),
	}
}

func (strg *MemStorage) Clear() {
	strg.mu.Lock()
	defer strg.mu.Unlock()
	strg.Counters = make(map[string]int64, 0)
	strg.Gauges = make(map[string]float64, 0)
}

func (strg *MemStorage) hasGauge(name string) bool {
	strg.mu.Lock()
	defer strg.mu.Unlock()
	_, ok := strg.Gauges[name]
	return ok
}

func (strg *MemStorage) hasCounter(name string) bool {
	strg.mu.Lock()
	defer strg.mu.Unlock()
	_, ok := strg.Counters[name]
	return ok
}

func (strg *MemStorage) AddCounter(ctx context.Context, key string, value int64) {
	strg.mu.Lock()
	defer strg.mu.Unlock()
	strg.Counters[key] += value
}

func (strg *MemStorage) SetGauge(ctx context.Context, key string, value float64) {
	strg.mu.Lock()
	defer strg.mu.Unlock()
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
	strg.mu.Lock()
	defer strg.mu.Unlock()
	val, ok := strg.Gauges[name]
	if !ok {
		return 0
	}
	return val
}

func (strg *MemStorage) GetCounterOrZero(ctx context.Context, name string) int64 {
	strg.mu.Lock()
	defer strg.mu.Unlock()
	val, ok := strg.Counters[name]
	if !ok {
		return 0
	}
	return val
}

func (strg *MemStorage) GetAllGauges(ctx context.Context) map[string]float64 {
	strg.mu.Lock()
	defer strg.mu.Unlock()

	return strg.Gauges
}

func (strg *MemStorage) GetAllCounters(ctx context.Context) map[string]int64 {
	strg.mu.Lock()
	defer strg.mu.Unlock()
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
