package sender

import (
	"github.com/gennadyterekhov/metrics-storage/internal/common/logger"

	"github.com/gennadyterekhov/metrics-storage/internal/agent/client"
	"github.com/gennadyterekhov/metrics-storage/internal/agent/metric"
)

type MetricsSender struct {
	Address         string
	IsGzip          bool
	Interval        int
	IsRunning       bool
	IsBatch         bool
	NumberOfWorkers int
}

func (msnd *MetricsSender) Report(memStatsPtr *metric.MetricsSet, metricsStorageClient *client.MetricsStorageClient) {
	msnd.IsRunning = true

	if msnd.IsBatch {
		msnd.sendAllMetricsInOneRequest(metricsStorageClient, memStatsPtr)
	} else {
		msnd.sendAllMetrics(metricsStorageClient, memStatsPtr)
	}

	msnd.IsRunning = false
}

func (msnd *MetricsSender) sendAllMetrics(metricsStorageClient *client.MetricsStorageClient, memStats *metric.MetricsSet) {
	jobs := make(chan metric.URLFormatter)

	for w := 1; w <= msnd.NumberOfWorkers; w++ {
		go worker(w, jobs, metricsStorageClient)
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

func (msnd *MetricsSender) sendAllMetricsInOneRequest(metricsStorageClient *client.MetricsStorageClient, memStats *metric.MetricsSet) {
	err := metricsStorageClient.SendAllMetricsInOneRequest(memStats)
	if err != nil {
		logger.Custom.Errorln("error when sending all metrics in one request ")
	}
}

func worker(workerIndex int, jobs <-chan metric.URLFormatter, metricsStorageClient *client.MetricsStorageClient) {
	for j := range jobs {
		logger.Custom.Debugln("worker", workerIndex, "started  job", j.GetName())
		err := metricsStorageClient.SendMetric(j)
		if err != nil {
			logger.Custom.Errorln("error when sending metric from worker ")
		}
		logger.Custom.Debugln("worker", workerIndex, "finished job", j.GetName())
	}
}
