package app

import (
	"github.com/gennadyterekhov/metrics-storage/internal/constants/types"
	"github.com/gennadyterekhov/metrics-storage/internal/domain/dto"
	"github.com/gennadyterekhov/metrics-storage/internal/logger"
	"github.com/gennadyterekhov/metrics-storage/internal/server/storage"
)

func SaveMetricToMemory(filledDto *dto.MetricToSaveDto) {
	logger.ZapSugarLogger.Debugln("saving metric",
		filledDto.Name, filledDto.Type, filledDto.CounterValue, filledDto.GaugeValue)
	if filledDto.Type == types.Counter {
		storage.MetricsRepository.AddCounter(filledDto.Name, filledDto.CounterValue)
	}
	if filledDto.Type == types.Gauge {
		storage.MetricsRepository.SetGauge(filledDto.Name, filledDto.GaugeValue)
	}
}
