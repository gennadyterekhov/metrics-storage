package requests

type GetMetricRequest struct {
	MetricType string `json:"type"`
	MetricName string `json:"id"`
	IsJSON     bool   `json:"-"`
	Error      error  `json:"-"`
}

type SaveMetricRequest struct {
	MetricType   string   `json:"type"`
	MetricName   string   `json:"id"`
	CounterValue *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	GaugeValue   *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
	IsJSON       bool     `json:"-"`
	Error        error    `json:"-"`
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

type SaveMetricListRequest []*SaveMetricRequest

func (r *SaveMetricRequest) GetMetricType() string   { return r.MetricType }
func (r *SaveMetricRequest) GetMetricName() string   { return r.MetricName }
func (r *SaveMetricRequest) GetCounterValue() *int64 { return r.CounterValue }
func (r *SaveMetricRequest) GetGaugeValue() *float64 { return r.GaugeValue }
func (r *SaveMetricRequest) GetIsJSON() bool         { return r.IsJSON }
