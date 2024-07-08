package services

import (
	"context"
	"fmt"

	"github.com/gennadyterekhov/metrics-storage/internal/server/repositories"

	"github.com/gennadyterekhov/metrics-storage/internal/common/constants/types"
	"github.com/gennadyterekhov/metrics-storage/internal/common/logger"
	"github.com/gennadyterekhov/metrics-storage/internal/server/httpui/requests"
	"github.com/gennadyterekhov/metrics-storage/internal/server/httpui/responses"
)

type GetMetricService struct {
	Repository repositories.RepositoryInterface
}

func NewGetMetricService(repo repositories.RepositoryInterface) GetMetricService {
	return GetMetricService{
		Repository: repo,
	}
}

func (srv GetMetricService) GetMetric(ctx context.Context, requestDto *requests.GetMetricRequest) (responseDto *responses.GetMetricResponse) {
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
		tmp, err := srv.Repository.GetCounter(ctx, requestDto.MetricName)
		responseDto.CounterValue, responseDto.Error = &tmp, err
	}
	if requestDto.MetricType == types.Gauge {
		tmp, err := srv.Repository.GetGauge(ctx, requestDto.MetricName)
		responseDto.GaugeValue, responseDto.Error = &tmp, err
	}
	return responseDto
}

func (srv GetMetricService) GetMetricsListAsHTML(ctx context.Context) string {
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
	gaugeList := srv.getGaugeList(ctx)
	counterList := srv.getCounterList(ctx)
	return fmt.Sprintf(
		templateText,
		types.Gauge,
		gaugeList,
		types.Counter,
		counterList,
	)
}

func (srv GetMetricService) getGaugeList(ctx context.Context) string {
	list := ""
	for name, val := range srv.Repository.GetAllGauges(ctx) {
		list += fmt.Sprintf("<li>%v : %v</li>", name, val)
	}

	return list
}

func (srv GetMetricService) getCounterList(ctx context.Context) string {
	list := ""
	for name, val := range srv.Repository.GetAllCounters(ctx) {
		list += fmt.Sprintf("<li>%v : %v</li>", name, val)
	}
	return list
}
