package services

import (
	"context"
	"time"

	"github.com/gennadyterekhov/metrics-storage/internal/server/repositories"

	"github.com/Rican7/retry"
	"github.com/Rican7/retry/backoff"
	"github.com/Rican7/retry/strategy"
	"github.com/gennadyterekhov/metrics-storage/internal/common/constants/types"
	"github.com/gennadyterekhov/metrics-storage/internal/common/logger"
	"github.com/gennadyterekhov/metrics-storage/internal/server/config"
	"github.com/gennadyterekhov/metrics-storage/internal/server/httpui/requests"
	"github.com/gennadyterekhov/metrics-storage/internal/server/httpui/responses"
)

type SaveMetricService struct {
	Repository repositories.RepositoryInterface
	Config     *config.ServerConfig
}

// ISaveMetric is used to decouple from requests.SaveMetricRequest
type ISaveMetric interface {
	GetMetricType() string
	GetMetricName() string
	GetCounterValue() *int64
	GetGaugeValue() *float64
	GetIsJSON() bool
}

func NewSaveMetricService(repo repositories.RepositoryInterface, conf *config.ServerConfig) SaveMetricService {
	return SaveMetricService{
		Repository: repo,
		Config:     conf,
	}
}

func (sms SaveMetricService) SaveMetricToMemory(ctx context.Context, filledDto ISaveMetric) (responseDto *responses.GetMetricResponse) {
	responseDto = &responses.GetMetricResponse{
		MetricType:   filledDto.GetMetricType(),
		MetricName:   filledDto.GetMetricName(),
		CounterValue: nil,
		GaugeValue:   nil,
		IsJSON:       filledDto.GetIsJSON(),
		Error:        nil,
	}
	logger.ZapSugarLogger.Debugln("saving metric",
		filledDto.GetMetricName(), filledDto.GetMetricType(), filledDto.GetCounterValue(), filledDto.GetGaugeValue())
	if filledDto.GetMetricType() == types.Counter && filledDto.GetCounterValue() != nil {
		sms.Repository.AddCounter(ctx, filledDto.GetMetricName(), *filledDto.GetCounterValue())
		updatedCounter, err := sms.Repository.GetCounter(ctx, filledDto.GetMetricName())
		if err != nil {
			updatedCounter = 0
		}
		responseDto.CounterValue = &updatedCounter
	}
	if filledDto.GetMetricType() == types.Gauge && filledDto.GetGaugeValue() != nil {
		sms.Repository.SetGauge(ctx, filledDto.GetMetricName(), *filledDto.GetGaugeValue())
		responseDto.GaugeValue = filledDto.GetGaugeValue()
	}

	sms.saveToDiskSynchronously(ctx)

	return responseDto
}

func (sms SaveMetricService) saveToDiskSynchronously(ctx context.Context) {
	if sms.Config.StoreInterval == 0 && sms.Config.FileStorage != "" {
		sms.SaveToDisk(ctx)
	}
}

func (sms SaveMetricService) SaveMetricListToMemory(ctx context.Context, filledDto *requests.SaveMetricListRequest) {
	logger.ZapSugarLogger.Debugln("saving metric list")
	for i := 0; i < len(*filledDto); i++ {
		sms.SaveMetricToMemory(ctx, (*filledDto)[i])
	}
}

func (sms SaveMetricService) SaveToDisk(ctx context.Context) {
	err := retry.Retry(
		func(attempt uint) error {
			return sms.Repository.SaveToDisk(ctx, sms.Config.FileStorage)
		},
		strategy.Limit(4),
		strategy.Backoff(backoff.Incremental(-1*time.Second, 2*time.Second)),
	)
	if err != nil {
		logger.ZapSugarLogger.Errorln("error when saving metric to file synchronously")
	}
}
