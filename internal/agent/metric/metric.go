package metric

import "fmt"

type URLFormatter interface {
	GetName() string
	GetType() string
	GetValueAsString() string
}

type CounterMetric struct {
	Name  string
	Type  string
	Value int64
}

type GaugeMetric struct {
	Name  string
	Type  string
	Value float64
}

type MetricsSet struct {
	Alloc         GaugeMetric
	BuckHashSys   GaugeMetric
	Frees         GaugeMetric
	GCCPUFraction GaugeMetric
	GCSys         GaugeMetric
	HeapAlloc     GaugeMetric
	HeapIdle      GaugeMetric
	HeapInuse     GaugeMetric
	HeapObjects   GaugeMetric
	HeapReleased  GaugeMetric
	HeapSys       GaugeMetric
	LastGC        GaugeMetric
	Lookups       GaugeMetric
	MCacheInuse   GaugeMetric
	MCacheSys     GaugeMetric
	MSpanInuse    GaugeMetric
	MSpanSys      GaugeMetric
	Mallocs       GaugeMetric
	NextGC        GaugeMetric
	NumForcedGC   GaugeMetric
	NumGC         GaugeMetric
	OtherSys      GaugeMetric
	PauseTotalNs  GaugeMetric
	StackInuse    GaugeMetric
	StackSys      GaugeMetric
	Sys           GaugeMetric
	TotalAlloc    GaugeMetric
	PollCount     CounterMetric
	RandomValue   GaugeMetric
	// so called additional metrics, introduced in 15 increment
	TotalMemory    GaugeMetric
	FreeMemory     GaugeMetric
	CPUUtilization []GaugeMetric
}

func (mtr *CounterMetric) GetName() string {
	return mtr.Name
}

func (mtr *CounterMetric) GetType() string {
	return mtr.Type
}

func (mtr *CounterMetric) GetValueAsString() string {
	return fmt.Sprintf("%v", mtr.Value)
}

func (mtr *GaugeMetric) GetName() string {
	return mtr.Name
}

func (mtr *GaugeMetric) GetType() string {
	return mtr.Type
}

func (mtr *GaugeMetric) GetValueAsString() string {
	return fmt.Sprintf("%v", mtr.Value)
}
