package repositories

type MetricsRepository interface {
	HasGauge(name string) bool
	HasCounter(name string) bool

	GetGauge(name string) float64
	GetCounter(name string) int64

	GetAllGauges() map[string]float64
	GetAllCounters() map[string]int64

	AddGauge(name string, value float64)
	AddCounter(name string, value int64)
}
