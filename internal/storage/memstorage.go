package storage

type MemStorage struct {
	Counters map[string]int64
	Gauges   map[string]float64
}

func CreateStorage() *MemStorage {
	return &MemStorage{
		make(map[string]int64, 0),
		make(map[string]float64, 0),
	}
}

func (strg *MemStorage) AddCounter(key string, value int64) {
	_, ok := strg.Counters[key]
	if ok {
		strg.Counters[key] += value
		return
	}
	strg.Counters[key] = value
}

func (strg *MemStorage) AddGauge(key string, value float64) {
	strg.Gauges[key] = value
}
