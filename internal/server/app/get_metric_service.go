package app

import (
	"fmt"
	"github.com/gennadyterekhov/metrics-storage/internal/constants/types"
	"github.com/gennadyterekhov/metrics-storage/internal/logger"
	"github.com/gennadyterekhov/metrics-storage/internal/server/httpui/requests"
	"github.com/gennadyterekhov/metrics-storage/internal/server/httpui/responses"
	"github.com/gennadyterekhov/metrics-storage/internal/server/storage"
)

func GetMetric(requestDto *requests.GetMetricRequest) (responseDto *responses.GetMetricResponse) {
	logger.ZapSugarLogger.Debugln("GetMetricAsString name, metricType", requestDto.MetricName, requestDto.MetricType)
	responseDto = &responses.GetMetricResponse{
		MetricType:   requestDto.MetricType,
		MetricName:   requestDto.MetricName,
		CounterValue: nil,
		GaugeValue:   nil,
		IsJSON:       requestDto.IsJSON,
		Error:        nil,
	}
	if requestDto.MetricType == types.Counter {
		tmp, err := storage.MetricsRepository.GetCounter(requestDto.MetricName)
		responseDto.CounterValue, responseDto.Error = &tmp, err
	}
	if requestDto.MetricType == types.Gauge {
		tmp, err := storage.MetricsRepository.GetGauge(requestDto.MetricName)
		responseDto.GaugeValue, responseDto.Error = &tmp, err
	}
	return responseDto
}

func GetMetricsListAsHTML() string {
	templateText := `
<!DOCTYPE html>
<html>
  <head></head>
  <body>
    <h2>%v</h2>
    <ul>
%v
    </ul>
    <h2>%v</h2>
    <ul>
%v
    </ul>
  </body>
</html>
`
	gaugeList := getGaugeList()
	counterList := getCounterList()
	return fmt.Sprintf(
		templateText,
		types.Gauge,
		gaugeList,
		types.Counter,
		counterList,
	)
}

func getGaugeList() string {
	list := ""
	for name, val := range storage.MetricsRepository.GetAllGauges() {
		list += fmt.Sprintf("<li>%v : %v</li>", name, val)
	}

	return list
}

func getCounterList() string {
	list := ""
	for name, val := range storage.MetricsRepository.GetAllCounters() {
		list += fmt.Sprintf("<li>%v : %v</li>", name, val)
	}
	return list
}
