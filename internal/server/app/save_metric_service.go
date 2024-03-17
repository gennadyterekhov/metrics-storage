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
		IsJSON:       filledDto.IsJSON,
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

	saveToDisk()

	return responseDto
}

func saveToDisk() {
	if config.Conf.StoreInterval == 0 && config.Conf.FileStorage != "" {
		err := storage.MetricsRepository.SaveToDisk(config.Conf.FileStorage)
		if err != nil {
			logger.ZapSugarLogger.Errorln("error when saving metric to file synchronously")
		}
	}
}

func SaveMetricBatchToMemory(filledDto *requests.SaveMetricBatchRequest) {
	// TODO refactor when db, can use fewer queries
	logger.ZapSugarLogger.Debugln("saving metric batch")

	setGaugeIfInDto(filledDto.Alloc)
	setGaugeIfInDto(filledDto.BuckHashSys)
	setGaugeIfInDto(filledDto.Frees)
	setGaugeIfInDto(filledDto.GCCPUFraction)
	setGaugeIfInDto(filledDto.GCSys)
	setGaugeIfInDto(filledDto.HeapAlloc)
	setGaugeIfInDto(filledDto.HeapIdle)
	setGaugeIfInDto(filledDto.HeapInuse)
	setGaugeIfInDto(filledDto.HeapObjects)
	setGaugeIfInDto(filledDto.HeapReleased)
	setGaugeIfInDto(filledDto.HeapSys)
	setGaugeIfInDto(filledDto.LastGC)
	setGaugeIfInDto(filledDto.Lookups)
	setGaugeIfInDto(filledDto.MCacheInuse)
	setGaugeIfInDto(filledDto.MCacheSys)
	setGaugeIfInDto(filledDto.MSpanInuse)
	setGaugeIfInDto(filledDto.MSpanSys)
	setGaugeIfInDto(filledDto.Mallocs)
	setGaugeIfInDto(filledDto.NextGC)
	setGaugeIfInDto(filledDto.NumForcedGC)
	setGaugeIfInDto(filledDto.NumGC)
	setGaugeIfInDto(filledDto.OtherSys)
	setGaugeIfInDto(filledDto.PauseTotalNs)
	setGaugeIfInDto(filledDto.StackInuse)
	setGaugeIfInDto(filledDto.StackSys)
	setGaugeIfInDto(filledDto.Sys)
	setGaugeIfInDto(filledDto.TotalAlloc)
	setGaugeIfInDto(filledDto.RandomValue)
	setCounterIfInDto(filledDto.PollCount)
	saveToDisk()
}

func SaveMetricListToMemory(filledDto *requests.SaveMetricListRequest) {
	logger.ZapSugarLogger.Debugln("saving metric list")
	for i := 0; i < len(*filledDto); i += 1 {
		SaveMetricToMemory(&(*filledDto)[i])
	}
}

func setGaugeIfInDto(filledDto *requests.GaugeMetricSubrequest) {
	if filledDto != nil {
		storage.MetricsRepository.SetGauge(filledDto.MetricName, filledDto.GaugeValue)
	}
}

func setCounterIfInDto(filledDto *requests.CounterMetricSubrequest) {
	if filledDto != nil {
		storage.MetricsRepository.AddCounter(filledDto.MetricName, filledDto.CounterValue)
	}
}
