package services

import (
	"context"
	"fmt"

	"github.com/gennadyterekhov/metrics-storage/internal/common/constants/types"
	"github.com/gennadyterekhov/metrics-storage/internal/common/logger"
	"github.com/gennadyterekhov/metrics-storage/internal/server/httpui/requests"
	"github.com/gennadyterekhov/metrics-storage/internal/server/httpui/responses"
	"github.com/gennadyterekhov/metrics-storage/internal/server/storage"
)

func GetMetric(ctx context.Context, requestDto *requests.GetMetricRequest) (responseDto *responses.GetMetricResponse) {
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
		tmp, err := storage.MetricsRepository.GetCounter(ctx, requestDto.MetricName)
		responseDto.CounterValue, responseDto.Error = &tmp, err
	}
	if requestDto.MetricType == types.Gauge {
		tmp, err := storage.MetricsRepository.GetGauge(ctx, requestDto.MetricName)
		responseDto.GaugeValue, responseDto.Error = &tmp, err
	}
	return responseDto
}

func GetMetricsListAsHTML(ctx context.Context) string {
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
	gaugeList := getGaugeList(ctx)
	counterList := getCounterList(ctx)
	return fmt.Sprintf(
		templateText,
		types.Gauge,
		gaugeList,
		types.Counter,
		counterList,
	)
}

func getGaugeList(ctx context.Context) string {
	list := ""
	for name, val := range storage.MetricsRepository.GetAllGauges(ctx) {
		list += fmt.Sprintf("<li>%v : %v</li>", name, val)
	}

	return list
}

func getCounterList(ctx context.Context) string {
	list := ""
	for name, val := range storage.MetricsRepository.GetAllCounters(ctx) {
		list += fmt.Sprintf("<li>%v : %v</li>", name, val)
	}
	return list
}
