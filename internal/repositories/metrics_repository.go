package repositories

type MetricsRepository interface {
	Clear()

	HasGauge(name string) bool
	HasCounter(name string) bool

	GetGauge(name string) (float64, error)
	GetCounter(name string) (int64, error)
	GetGaugeOrZero(name string) float64
	GetCounterOrZero(name string) int64

	GetAllGauges() map[string]float64
	GetAllCounters() map[string]int64

	AddGauge(name string, value float64)
	AddCounter(name string, value int64)
}
