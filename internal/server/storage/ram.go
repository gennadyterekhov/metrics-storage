package storage

import (
	"fmt"
	"github.com/gennadyterekhov/metrics-storage/internal/constants/exceptions"
)

type MemStorage struct {
	Counters map[string]int64
	Gauges   map[string]float64
}

func CreateStorage() *MemStorage {
	return &MemStorage{
		Counters: make(map[string]int64, 0),
		Gauges:   make(map[string]float64, 0),
	}
}

func (strg *MemStorage) Clear() {
	strg.Counters = make(map[string]int64, 0)
	strg.Gauges = make(map[string]float64, 0)
}

func (strg *MemStorage) HasGauge(name string) bool {
	_, ok := strg.Gauges[name]
	return ok
}

func (strg *MemStorage) HasCounter(name string) bool {
	_, ok := strg.Counters[name]
	return ok
}

func (strg *MemStorage) AddCounter(key string, value int64) {
	strg.Counters[key] += value
}

func (strg *MemStorage) SetGauge(key string, value float64) {
	strg.Gauges[key] = value
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
	val, ok := strg.Gauges[name]
	if !ok {
		return 0
	}
	return val
}

func (strg *MemStorage) GetCounterOrZero(name string) int64 {
	val, ok := strg.Counters[name]
	if !ok {
		return 0
	}
	return val
}

func (strg *MemStorage) GetAllGauges() map[string]float64 {
	return strg.Gauges
}

func (strg *MemStorage) GetAllCounters() map[string]int64 {
	return strg.Counters
}

func (strg *MemStorage) GetAll() (map[string]float64, map[string]int64) {
	return strg.GetAllGauges(), strg.GetAllCounters()
}
