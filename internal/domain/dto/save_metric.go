package dto

/*
TODO merge with metrics.Metrics
*/
type MetricToSaveDto struct {
	Type         string
	Name         string
	CounterValue int64
	GaugeValue   float64
}
