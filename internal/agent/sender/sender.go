package sender

import (
	"github.com/gennadyterekhov/metrics-storage/internal/agent/client"
	"github.com/gennadyterekhov/metrics-storage/internal/agent/metric"
	"github.com/gennadyterekhov/metrics-storage/internal/logger"
	"time"
)

type MetricsSender struct {
	Address   string
	Interval  int
	IsRunning bool
}

func (msnd *MetricsSender) wait() {
	time.Sleep(time.Duration(msnd.Interval * int(time.Second)))
}

func (msnd *MetricsSender) Report(memStatsPtr *metric.MetricsSet) {
	msnd.IsRunning = true
	msnd.wait()

	sendAllMetrics(msnd.Address, memStatsPtr)
	msnd.IsRunning = false
}

func sendAllMetrics(address string, memStats *metric.MetricsSet) {
	logger.ZapSugarLogger.Debugln("gonna send metrics to server")

	sendMetric(&memStats.Alloc, address)
	sendMetric(&memStats.BuckHashSys, address)
	sendMetric(&memStats.Frees, address)
	sendMetric(&memStats.GCCPUFraction, address)
	sendMetric(&memStats.GCSys, address)
	sendMetric(&memStats.HeapAlloc, address)
	sendMetric(&memStats.HeapIdle, address)
	sendMetric(&memStats.HeapInuse, address)
	sendMetric(&memStats.HeapObjects, address)
	sendMetric(&memStats.HeapReleased, address)
	sendMetric(&memStats.HeapSys, address)
	sendMetric(&memStats.LastGC, address)
	sendMetric(&memStats.Lookups, address)
	sendMetric(&memStats.MCacheInuse, address)
	sendMetric(&memStats.MCacheSys, address)
	sendMetric(&memStats.MSpanInuse, address)
	sendMetric(&memStats.MSpanSys, address)
	sendMetric(&memStats.Mallocs, address)
	sendMetric(&memStats.NextGC, address)
	sendMetric(&memStats.NumForcedGC, address)
	sendMetric(&memStats.NumGC, address)
	sendMetric(&memStats.OtherSys, address)
	sendMetric(&memStats.PauseTotalNs, address)
	sendMetric(&memStats.StackInuse, address)
	sendMetric(&memStats.StackSys, address)
	sendMetric(&memStats.Sys, address)
	sendMetric(&memStats.TotalAlloc, address)
	sendMetric(&memStats.PollCount, address)
	sendMetric(&memStats.RandomValue, address)
}

func sendMetric(met metric.MetricURLFormatter, address string) error {
	err := client.SendMetric(met, address)

	if err != nil {
		logger.ZapSugarLogger.Errorln("error when sending metric to server", err.Error())

		return err
	}
	return nil
}
