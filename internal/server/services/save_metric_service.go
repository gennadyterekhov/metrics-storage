package services

import (
	"context"
	"time"

	"github.com/Rican7/retry"
	"github.com/Rican7/retry/backoff"
	"github.com/Rican7/retry/strategy"
	"github.com/gennadyterekhov/metrics-storage/internal/common/constants/types"
	"github.com/gennadyterekhov/metrics-storage/internal/common/logger"
	"github.com/gennadyterekhov/metrics-storage/internal/server/config"
	"github.com/gennadyterekhov/metrics-storage/internal/server/httpui/requests"
	"github.com/gennadyterekhov/metrics-storage/internal/server/httpui/responses"
	"github.com/gennadyterekhov/metrics-storage/internal/server/storage"
)

func SaveMetricToMemory(ctx context.Context, filledDto *requests.SaveMetricRequest) (responseDto *responses.GetMetricResponse) {
	responseDto = &responses.GetMetricResponse{
		MetricType:   filledDto.MetricType,
		MetricName:   filledDto.MetricName,
		CounterValue: nil,
		GaugeValue:   nil,
		IsJSON:       filledDto.IsJSON,
		Error:        nil,
	}
	logger.ZapSugarLogger.Debugln("saving metric",
		filledDto.MetricName, filledDto.MetricType, filledDto.CounterValue, filledDto.GaugeValue)
	if filledDto.MetricType == types.Counter && filledDto.CounterValue != nil {
		storage.MetricsRepository.AddCounter(ctx, filledDto.MetricName, *filledDto.CounterValue)
		updatedCounter := storage.MetricsRepository.GetCounterOrZero(ctx, filledDto.MetricName)
		responseDto.CounterValue = &updatedCounter
	}
	if filledDto.MetricType == types.Gauge && filledDto.GaugeValue != nil {
		storage.MetricsRepository.SetGauge(ctx, filledDto.MetricName, *filledDto.GaugeValue)
		responseDto.GaugeValue = filledDto.GaugeValue
	}

	saveToDiskSynchronously(ctx)

	return responseDto
}

func saveToDiskSynchronously(ctx context.Context) {
	if config.Conf.StoreInterval == 0 && config.Conf.FileStorage != "" {
		SaveToDisk(ctx)
	}
}

func SaveToDisk(ctx context.Context) {
	err := retry.Retry(
		func(attempt uint) error {
			return storage.MetricsRepository.SaveToDisk(ctx, config.Conf.FileStorage)
		},
		strategy.Limit(4),
		strategy.Backoff(backoff.Incremental(-1*time.Second, 2*time.Second)),
	)
	if err != nil {
		logger.ZapSugarLogger.Errorln("error when saving metric to file synchronously")
	}
}

func SaveMetricListToMemory(ctx context.Context, filledDto *requests.SaveMetricListRequest) {
	logger.ZapSugarLogger.Debugln("saving metric list")
	for i := 0; i < len(*filledDto); i += 1 {
		SaveMetricToMemory(ctx, (*filledDto)[i])
	}
}
