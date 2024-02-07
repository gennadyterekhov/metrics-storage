package repositories

type MetricsRepository interface {
	GetGauge(name string) (float64, error)
	GetCounter(name string) (int64, error)

	GetAllGauges() map[string]float64
	GetAllCounters() map[string]int64
	GetAll() (map[string]float64, map[string]int64)

	SetGauge(name string, value float64)
	AddCounter(name string, value int64)
}
