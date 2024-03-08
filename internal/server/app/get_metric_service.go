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
		CounterValue: 0,
		GaugeValue:   0,
		IsJson:       requestDto.IsJson,
		Error:        nil,
	}
	if requestDto.MetricType == types.Counter {
		responseDto.CounterValue, responseDto.Error = storage.MetricsRepository.GetCounter(requestDto.MetricName)
	}
	if requestDto.MetricType == types.Gauge {
		responseDto.GaugeValue, responseDto.Error = storage.MetricsRepository.GetGauge(requestDto.MetricName)
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
