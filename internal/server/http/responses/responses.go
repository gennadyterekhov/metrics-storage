package responses

type GetMetricResponse struct {
	MetricType   string   `json:"type"`
	MetricName   string   `json:"id"`
	CounterValue *int64   `json:"delta,omitempty"`
	GaugeValue   *float64 `json:"value,omitempty"`
	IsJSON       bool     `json:"-"`
	Error        error    `json:"-"`
}
