package requests

type GetMetricRequest struct {
	MetricType string `json:"type"`
	MetricName string `json:"id"`
	IsJson     bool   `json:"-"`
	Error      error  `json:"-"`
}

type SaveMetricRequest struct {
	MetricType string `json:"type"`
	MetricName string `json:"id"`

	CounterValue *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	GaugeValue   *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge

	IsJson bool  `json:"-"`
	Error  error `json:"-"`
}

type GaugeMetricSubrequest struct {
	MetricName string  `json:"id"`
	MetricType string  `json:"type"`
	GaugeValue float64 `json:"value"`
}

type CounterMetricSubrequest struct {
	MetricName   string `json:"id"`
	MetricType   string `json:"type"`
	CounterValue int64  `json:"delta"`
}

type SaveMetricBatchRequest struct {
	Error error `json:"-"`

	Alloc         GaugeMetricSubrequest   `json:"Alloc"`
	BuckHashSys   GaugeMetricSubrequest   `json:"BuckHashSys"`
	Frees         GaugeMetricSubrequest   `json:"Frees"`
	GCCPUFraction GaugeMetricSubrequest   `json:"GCCPUFraction"`
	GCSys         GaugeMetricSubrequest   `json:"GCSys"`
	HeapAlloc     GaugeMetricSubrequest   `json:"HeapAlloc"`
	HeapIdle      GaugeMetricSubrequest   `json:"HeapIdle"`
	HeapInuse     GaugeMetricSubrequest   `json:"HeapInuse"`
	HeapObjects   GaugeMetricSubrequest   `json:"HeapObjects"`
	HeapReleased  GaugeMetricSubrequest   `json:"HeapReleased"`
	HeapSys       GaugeMetricSubrequest   `json:"HeapSys"`
	LastGC        GaugeMetricSubrequest   `json:"LastGC"`
	Lookups       GaugeMetricSubrequest   `json:"Lookups"`
	MCacheInuse   GaugeMetricSubrequest   `json:"MCacheInuse"`
	MCacheSys     GaugeMetricSubrequest   `json:"MCacheSys"`
	MSpanInuse    GaugeMetricSubrequest   `json:"MSpanInuse"`
	MSpanSys      GaugeMetricSubrequest   `json:"MSpanSys"`
	Mallocs       GaugeMetricSubrequest   `json:"Mallocs"`
	NextGC        GaugeMetricSubrequest   `json:"NextGC"`
	NumForcedGC   GaugeMetricSubrequest   `json:"NumForcedGC"`
	NumGC         GaugeMetricSubrequest   `json:"NumGC"`
	OtherSys      GaugeMetricSubrequest   `json:"OtherSys"`
	PauseTotalNs  GaugeMetricSubrequest   `json:"PauseTotalNs"`
	StackInuse    GaugeMetricSubrequest   `json:"StackInuse"`
	StackSys      GaugeMetricSubrequest   `json:"StackSys"`
	Sys           GaugeMetricSubrequest   `json:"Sys"`
	TotalAlloc    GaugeMetricSubrequest   `json:"TotalAlloc"`
	PollCount     CounterMetricSubrequest `json:"PollCount"`
	RandomValue   GaugeMetricSubrequest   `json:"RandomValue"`
}
