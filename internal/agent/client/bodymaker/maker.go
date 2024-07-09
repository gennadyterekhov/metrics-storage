package bodymaker

import (
	"encoding/json"
	"strconv"

	"github.com/gennadyterekhov/metrics-storage/internal/agent/metric"
	"github.com/gennadyterekhov/metrics-storage/internal/common/constants/types"
	"github.com/gennadyterekhov/metrics-storage/internal/common/logger"
	"github.com/gennadyterekhov/metrics-storage/internal/server/httpui/requests"
)

type MetricsRequest struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func GetBody(met metric.MetricURLFormatter) ([]byte, error) {
	counterValue, gaugeValue, err := getMetricValues(met)
	if err != nil {
		return nil, err
	}

	metricToEncode := MetricsRequest{
		ID:    met.GetName(),
		MType: met.GetType(),
		Delta: &counterValue,
		Value: &gaugeValue,
	}
	jsonBytes, err := json.Marshal(metricToEncode)
	if err != nil {
		logger.ZapSugarLogger.Errorln("error when encoding metric", err.Error())

		return nil, err
	}
	return jsonBytes, nil
}

func GetBodyForAllMetrics(memStats *metric.MetricsSet) ([]byte, error) {
	metricToEncode := requests.SaveMetricListRequest{
		getSubrequest(&memStats.Alloc),
		getSubrequest(&memStats.BuckHashSys),
		getSubrequest(&memStats.Frees),
		getSubrequest(&memStats.GCCPUFraction),
		getSubrequest(&memStats.GCSys),
		getSubrequest(&memStats.HeapAlloc),
		getSubrequest(&memStats.HeapIdle),
		getSubrequest(&memStats.HeapInuse),
		getSubrequest(&memStats.HeapObjects),
		getSubrequest(&memStats.HeapReleased),
		getSubrequest(&memStats.HeapSys),
		getSubrequest(&memStats.LastGC),
		getSubrequest(&memStats.Lookups),
		getSubrequest(&memStats.MCacheInuse),
		getSubrequest(&memStats.MCacheSys),
		getSubrequest(&memStats.MSpanInuse),
		getSubrequest(&memStats.MSpanSys),
		getSubrequest(&memStats.Mallocs),
		getSubrequest(&memStats.NextGC),
		getSubrequest(&memStats.NumForcedGC),
		getSubrequest(&memStats.NumGC),
		getSubrequest(&memStats.OtherSys),
		getSubrequest(&memStats.PauseTotalNs),
		getSubrequest(&memStats.StackInuse),
		getSubrequest(&memStats.StackSys),
		getSubrequest(&memStats.Sys),
		getSubrequest(&memStats.TotalAlloc),
		getSubrequest(&memStats.PollCount),
		getSubrequest(&memStats.RandomValue),

		getSubrequest(&memStats.TotalMemory),
		getSubrequest(&memStats.FreeMemory),
	}

	for i := 0; i < len(memStats.CPUUtilization); i += 1 {
		metricToEncode = append(metricToEncode, getSubrequest(&memStats.CPUUtilization[i]))
	}

	jsonBytes, err := json.Marshal(metricToEncode)
	if err != nil {
		logger.ZapSugarLogger.Errorln("error when encoding metric batch", err.Error())

		return nil, err
	}
	return jsonBytes, nil
}

func getSubrequest(met metric.MetricURLFormatter) *requests.SaveMetricRequest {
	counter, gauge, err := getMetricValues(met)
	if err != nil {
		return nil
	}

	return &requests.SaveMetricRequest{
		MetricName:   met.GetName(),
		MetricType:   met.GetType(),
		CounterValue: &counter,
		GaugeValue:   &gauge,
	}
}

func getMetricValues(met metric.MetricURLFormatter) (counterValue int64, gaugeValue float64, err error) {
	if met.GetType() == types.Counter {
		counterValue, err = strconv.ParseInt(met.GetValueAsString(), 10, 64)
		if err != nil {
			return 0, 0, err
		}
	} else {
		gaugeValue, err = strconv.ParseFloat(met.GetValueAsString(), 64)
		if err != nil {
			return 0, 0, err
		}
	}

	return counterValue, gaugeValue, nil
}
