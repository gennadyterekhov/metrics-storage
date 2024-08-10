package sender

import (
	"github.com/gennadyterekhov/metrics-storage/internal/agent/config"
	"github.com/gennadyterekhov/metrics-storage/internal/common/logger"

	"github.com/gennadyterekhov/metrics-storage/internal/agent/client"
	"github.com/gennadyterekhov/metrics-storage/internal/agent/metric"
)

type MetricsSender struct {
	MetricsStorageClient *client.MetricsStorageClient
	Address              string
	IsGzip               bool
	Interval             int
	IsRunning            bool
	IsBatch              bool
	NumberOfWorkers      int
}

func New(metricsStorageClient *client.MetricsStorageClient, conf *config.Config) *MetricsSender {
	return &MetricsSender{
		MetricsStorageClient: metricsStorageClient,
		Address:              conf.Addr,
		IsGzip:               conf.IsGzip,
		Interval:             conf.ReportInterval,
		IsRunning:            false,
		IsBatch:              conf.IsBatch,
		NumberOfWorkers:      conf.SimultaneousRequestsLimit,
	}
}

func (msnd *MetricsSender) Report(memStatsPtr *metric.MetricsSet) {
	msnd.IsRunning = true

	if msnd.IsBatch {
		msnd.sendAllMetricsInOneRequest(memStatsPtr)
	} else {
		msnd.sendAllMetrics(memStatsPtr)
	}

	msnd.IsRunning = false
}

func (msnd *MetricsSender) sendAllMetrics(memStats *metric.MetricsSet) {
	jobs := make(chan metric.URLFormatter)

	for w := 1; w <= msnd.NumberOfWorkers; w++ {
		go msnd.worker(w, jobs)
	}
	setJobs(jobs, memStats)
	close(jobs)
}

func setJobs(jobs chan metric.URLFormatter, memStats *metric.MetricsSet) {
	jobs <- &memStats.Alloc
	jobs <- &memStats.BuckHashSys
	jobs <- &memStats.Frees
	jobs <- &memStats.GCCPUFraction
	jobs <- &memStats.GCSys
	jobs <- &memStats.HeapAlloc
	jobs <- &memStats.HeapIdle
	jobs <- &memStats.HeapInuse
	jobs <- &memStats.HeapObjects
	jobs <- &memStats.HeapReleased
	jobs <- &memStats.HeapSys
	jobs <- &memStats.LastGC
	jobs <- &memStats.Lookups
	jobs <- &memStats.MCacheInuse
	jobs <- &memStats.MCacheSys
	jobs <- &memStats.MSpanInuse
	jobs <- &memStats.MSpanSys
	jobs <- &memStats.Mallocs
	jobs <- &memStats.NextGC
	jobs <- &memStats.NumForcedGC
	jobs <- &memStats.NumGC
	jobs <- &memStats.OtherSys
	jobs <- &memStats.PauseTotalNs
	jobs <- &memStats.StackInuse
	jobs <- &memStats.StackSys
	jobs <- &memStats.Sys
	jobs <- &memStats.TotalAlloc
	jobs <- &memStats.PollCount
	jobs <- &memStats.RandomValue

	jobs <- &memStats.TotalMemory
	jobs <- &memStats.FreeMemory
	for i := 0; i < len(memStats.CPUUtilization); i++ {
		jobs <- &memStats.CPUUtilization[i]
	}
}

func (msnd *MetricsSender) sendAllMetricsInOneRequest(memStats *metric.MetricsSet) {
	err := msnd.MetricsStorageClient.SendAllMetricsInOneRequest(memStats)
	if err != nil {
		logger.Custom.Errorln("error when sending all metrics in one request ")
	}
}

func (msnd *MetricsSender) worker(workerIndex int, jobs <-chan metric.URLFormatter) {
	for j := range jobs {
		logger.Custom.Debugln("worker", workerIndex, "started  job", j.GetName())
		err := msnd.MetricsStorageClient.SendMetric(j)
		if err != nil {
			logger.Custom.Errorln("error when sending metric from worker ")
		}
		logger.Custom.Debugln("worker", workerIndex, "finished job", j.GetName())
	}
}
