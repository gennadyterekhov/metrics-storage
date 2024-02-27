package savemetricservice

import (
	"github.com/gennadyterekhov/metrics-storage/internal/constants/types"
	"github.com/gennadyterekhov/metrics-storage/internal/container"
	"github.com/gennadyterekhov/metrics-storage/internal/domain/dto"
	"github.com/gennadyterekhov/metrics-storage/internal/logger"
)

func SaveMetricToMemory(filledDto *dto.MetricToSaveDto) {
	logger.ZapSugarLogger.Debugln("saving metric",
		filledDto.Name, filledDto.Type, filledDto.CounterValue, filledDto.GaugeValue)
	if filledDto.Type == types.Counter {
		container.MetricsRepository.AddCounter(filledDto.Name, filledDto.CounterValue)
	}
	if filledDto.Type == types.Gauge {
		container.MetricsRepository.SetGauge(filledDto.Name, filledDto.GaugeValue)
	}
}
