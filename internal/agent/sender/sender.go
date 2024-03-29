package sender

import (
	"github.com/gennadyterekhov/metrics-storage/internal/agent/client"
	"github.com/gennadyterekhov/metrics-storage/internal/agent/metric"
)

type MetricsSender struct {
	Address   string
	IsGzip    bool
	Interval  int
	IsRunning bool
	IsBatch   bool
}

func (msnd *MetricsSender) Report(memStatsPtr *metric.MetricsSet, metricsStorageClient *client.MetricsStorageClient) {
	msnd.IsRunning = true

	if msnd.IsBatch {
		sendAllMetricsInOneRequest(metricsStorageClient, memStatsPtr)
	} else {
		sendAllMetrics(metricsStorageClient, memStatsPtr)
	}

	msnd.IsRunning = false
}

func sendAllMetrics(metricsStorageClient *client.MetricsStorageClient, memStats *metric.MetricsSet) {
	metricsStorageClient.SendMetric(&memStats.Alloc)
	metricsStorageClient.SendMetric(&memStats.BuckHashSys)
	metricsStorageClient.SendMetric(&memStats.Frees)
	metricsStorageClient.SendMetric(&memStats.GCCPUFraction)
	metricsStorageClient.SendMetric(&memStats.GCSys)
	metricsStorageClient.SendMetric(&memStats.HeapAlloc)
	metricsStorageClient.SendMetric(&memStats.HeapIdle)
	metricsStorageClient.SendMetric(&memStats.HeapInuse)
	metricsStorageClient.SendMetric(&memStats.HeapObjects)
	metricsStorageClient.SendMetric(&memStats.HeapReleased)
	metricsStorageClient.SendMetric(&memStats.HeapSys)
	metricsStorageClient.SendMetric(&memStats.LastGC)
	metricsStorageClient.SendMetric(&memStats.Lookups)
	metricsStorageClient.SendMetric(&memStats.MCacheInuse)
	metricsStorageClient.SendMetric(&memStats.MCacheSys)
	metricsStorageClient.SendMetric(&memStats.MSpanInuse)
	metricsStorageClient.SendMetric(&memStats.MSpanSys)
	metricsStorageClient.SendMetric(&memStats.Mallocs)
	metricsStorageClient.SendMetric(&memStats.NextGC)
	metricsStorageClient.SendMetric(&memStats.NumForcedGC)
	metricsStorageClient.SendMetric(&memStats.NumGC)
	metricsStorageClient.SendMetric(&memStats.OtherSys)
	metricsStorageClient.SendMetric(&memStats.PauseTotalNs)
	metricsStorageClient.SendMetric(&memStats.StackInuse)
	metricsStorageClient.SendMetric(&memStats.StackSys)
	metricsStorageClient.SendMetric(&memStats.Sys)
	metricsStorageClient.SendMetric(&memStats.TotalAlloc)
	metricsStorageClient.SendMetric(&memStats.PollCount)
	metricsStorageClient.SendMetric(&memStats.RandomValue)
}

func sendAllMetricsInOneRequest(metricsStorageClient *client.MetricsStorageClient, memStats *metric.MetricsSet) {
	metricsStorageClient.SendAllMetricsInOneRequest(memStats)
}
