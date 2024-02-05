package storage

import (
	"fmt"
	"github.com/gennadyterekhov/metrics-storage/internal/exceptions"
	"github.com/gennadyterekhov/metrics-storage/internal/repositories"
)

type MemStorage struct {
	counters map[string]int64
	gauges   map[string]float64
}

func CreateStorage() repositories.MetricsRepository {
	return &MemStorage{
		counters: make(map[string]int64, 0),
		gauges:   make(map[string]float64, 0),
	}
}

func (strg *MemStorage) Clear() {
	strg.counters = make(map[string]int64, 0)
	strg.gauges = make(map[string]float64, 0)
}

func (strg *MemStorage) HasGauge(name string) bool {
	_, ok := strg.gauges[name]
	return ok
}

func (strg *MemStorage) HasCounter(name string) bool {
	_, ok := strg.counters[name]
	return ok
}

func (strg *MemStorage) AddCounter(key string, value int64) {
	_, ok := strg.counters[key]
	if ok {
		strg.counters[key] += value
		return
	}
	strg.counters[key] = value
}

func (strg *MemStorage) AddGauge(key string, value float64) {
	strg.gauges[key] = value
}

func (strg *MemStorage) GetGauge(name string) (float64, error) {
	if !strg.HasGauge(name) {
		return 0, fmt.Errorf(exceptions.UnknownMetricName)
	}
	return strg.GetGaugeOrZero(name), nil
}

func (strg *MemStorage) GetCounter(name string) (int64, error) {
	if !strg.HasCounter(name) {
		return 0, fmt.Errorf(exceptions.UnknownMetricName)
	}
	return strg.GetCounterOrZero(name), nil
}

func (strg *MemStorage) GetGaugeOrZero(name string) float64 {
	val, ok := strg.gauges[name]
	if !ok {
		return 0
	}
	return val
}

func (strg *MemStorage) GetCounterOrZero(name string) int64 {
	val, ok := strg.counters[name]
	if !ok {
		return 0
	}
	return val
}

func (strg *MemStorage) GetAllGauges() map[string]float64 {
	return strg.gauges
}

func (strg *MemStorage) GetAllCounters() map[string]int64 {
	return strg.counters
}
