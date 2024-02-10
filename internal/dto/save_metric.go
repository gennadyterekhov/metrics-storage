package dto

type MetricToSaveDto struct {
	Type         string
	Name         string
	CounterValue int64
	GaugeValue   float64
}
