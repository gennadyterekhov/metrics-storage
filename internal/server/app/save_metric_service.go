package app

import (
	"github.com/gennadyterekhov/metrics-storage/internal/constants/types"
	"github.com/gennadyterekhov/metrics-storage/internal/logger"
	"github.com/gennadyterekhov/metrics-storage/internal/server/config"
	"github.com/gennadyterekhov/metrics-storage/internal/server/httpui/requests"
	"github.com/gennadyterekhov/metrics-storage/internal/server/httpui/responses"
	"github.com/gennadyterekhov/metrics-storage/internal/server/storage"
)

func SaveMetricToMemory(filledDto *requests.SaveMetricRequest) (responseDto *responses.GetMetricResponse) {
	responseDto = &responses.GetMetricResponse{
		MetricType:   filledDto.MetricType,
		MetricName:   filledDto.MetricName,
		CounterValue: nil,
		GaugeValue:   nil,
		IsJson:       filledDto.IsJson,
		Error:        nil,
	}
	logger.ZapSugarLogger.Debugln("saving metric",
		filledDto.MetricName, filledDto.MetricType, filledDto.CounterValue, filledDto.GaugeValue)
	if filledDto.MetricType == types.Counter && filledDto.CounterValue != nil {
		storage.MetricsRepository.AddCounter(filledDto.MetricName, *filledDto.CounterValue)
		updatedCounter := storage.MetricsRepository.GetCounterOrZero(filledDto.MetricName)
		responseDto.CounterValue = &updatedCounter
	}
	if filledDto.MetricType == types.Gauge && filledDto.GaugeValue != nil {
		storage.MetricsRepository.SetGauge(filledDto.MetricName, *filledDto.GaugeValue)
		responseDto.GaugeValue = filledDto.GaugeValue
	}

	if config.Conf.StoreInterval == 0 && config.Conf.FileStorage != "" {
		err := storage.MetricsRepository.Save(config.Conf.FileStorage)
		if err != nil {
			logger.ZapSugarLogger.Errorln("error when saving metric to file synchronously")
		}
	}

	return responseDto
}

func SaveMetricBatchToMemory(filledDto *requests.SaveMetricBatchRequest) {
	logger.ZapSugarLogger.Debugln("saving metric batch")

	storage.MetricsRepository.SetGauge(filledDto.Alloc.MetricName, filledDto.Alloc.GaugeValue)
	storage.MetricsRepository.SetGauge(filledDto.BuckHashSys.MetricName, filledDto.BuckHashSys.GaugeValue)
	storage.MetricsRepository.SetGauge(filledDto.Frees.MetricName, filledDto.Frees.GaugeValue)
	storage.MetricsRepository.SetGauge(filledDto.GCCPUFraction.MetricName, filledDto.GCCPUFraction.GaugeValue)
	storage.MetricsRepository.SetGauge(filledDto.GCSys.MetricName, filledDto.GCSys.GaugeValue)
	storage.MetricsRepository.SetGauge(filledDto.HeapAlloc.MetricName, filledDto.HeapAlloc.GaugeValue)
	storage.MetricsRepository.SetGauge(filledDto.HeapIdle.MetricName, filledDto.HeapIdle.GaugeValue)
	storage.MetricsRepository.SetGauge(filledDto.HeapInuse.MetricName, filledDto.HeapInuse.GaugeValue)
	storage.MetricsRepository.SetGauge(filledDto.HeapObjects.MetricName, filledDto.HeapObjects.GaugeValue)
	storage.MetricsRepository.SetGauge(filledDto.HeapReleased.MetricName, filledDto.HeapReleased.GaugeValue)
	storage.MetricsRepository.SetGauge(filledDto.HeapSys.MetricName, filledDto.HeapSys.GaugeValue)
	storage.MetricsRepository.SetGauge(filledDto.LastGC.MetricName, filledDto.LastGC.GaugeValue)
	storage.MetricsRepository.SetGauge(filledDto.Lookups.MetricName, filledDto.Lookups.GaugeValue)
	storage.MetricsRepository.SetGauge(filledDto.MCacheInuse.MetricName, filledDto.MCacheInuse.GaugeValue)
	storage.MetricsRepository.SetGauge(filledDto.MCacheSys.MetricName, filledDto.MCacheSys.GaugeValue)
	storage.MetricsRepository.SetGauge(filledDto.MSpanInuse.MetricName, filledDto.MSpanInuse.GaugeValue)
	storage.MetricsRepository.SetGauge(filledDto.MSpanSys.MetricName, filledDto.MSpanSys.GaugeValue)
	storage.MetricsRepository.SetGauge(filledDto.Mallocs.MetricName, filledDto.Mallocs.GaugeValue)
	storage.MetricsRepository.SetGauge(filledDto.NextGC.MetricName, filledDto.NextGC.GaugeValue)
	storage.MetricsRepository.SetGauge(filledDto.NumForcedGC.MetricName, filledDto.NumForcedGC.GaugeValue)
	storage.MetricsRepository.SetGauge(filledDto.NumGC.MetricName, filledDto.NumGC.GaugeValue)
	storage.MetricsRepository.SetGauge(filledDto.OtherSys.MetricName, filledDto.OtherSys.GaugeValue)
	storage.MetricsRepository.SetGauge(filledDto.PauseTotalNs.MetricName, filledDto.PauseTotalNs.GaugeValue)
	storage.MetricsRepository.SetGauge(filledDto.StackInuse.MetricName, filledDto.StackInuse.GaugeValue)
	storage.MetricsRepository.SetGauge(filledDto.StackSys.MetricName, filledDto.StackSys.GaugeValue)
	storage.MetricsRepository.SetGauge(filledDto.Sys.MetricName, filledDto.Sys.GaugeValue)
	storage.MetricsRepository.SetGauge(filledDto.TotalAlloc.MetricName, filledDto.TotalAlloc.GaugeValue)
	storage.MetricsRepository.SetGauge(filledDto.RandomValue.MetricName, filledDto.RandomValue.GaugeValue)
	storage.MetricsRepository.AddCounter(filledDto.PollCount.MetricName, filledDto.PollCount.CounterValue)
}
