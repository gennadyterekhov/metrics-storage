package sender

import (
	"fmt"
	"github.com/gennadyterekhov/metrics-storage/internal/agent/metric"
	"github.com/go-resty/resty/v2"
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
	urls := getURLs(memStats)
	for i := 0; i < len(urls); i++ {
		_ = sendMetric(address + urls[i])
	}
}

func getURLs(memStats *metric.MetricsSet) []string {
	return []string{
		getURL(&memStats.Alloc),
		getURL(&memStats.BuckHashSys),
		getURL(&memStats.Frees),
		getURL(&memStats.GCCPUFraction),
		getURL(&memStats.GCSys),
		getURL(&memStats.HeapAlloc),
		getURL(&memStats.HeapIdle),
		getURL(&memStats.HeapInuse),
		getURL(&memStats.HeapObjects),
		getURL(&memStats.HeapReleased),
		getURL(&memStats.HeapSys),
		getURL(&memStats.LastGC),
		getURL(&memStats.Lookups),
		getURL(&memStats.MCacheInuse),
		getURL(&memStats.MCacheSys),
		getURL(&memStats.MSpanInuse),
		getURL(&memStats.MSpanSys),
		getURL(&memStats.Mallocs),
		getURL(&memStats.NextGC),
		getURL(&memStats.NumForcedGC),
		getURL(&memStats.NumGC),
		getURL(&memStats.OtherSys),
		getURL(&memStats.PauseTotalNs),
		getURL(&memStats.StackInuse),
		getURL(&memStats.StackSys),
		getURL(&memStats.Sys),
		getURL(&memStats.TotalAlloc),
	}
}

func getURL(met metric.MerticURLFormatter) string {
	template := "/update/%v/%v/%v"
	return fmt.Sprintf(template, met.GetType(), met.GetName(), met.GetValueAsString())
}

func sendMetric(url string) (err error) {
	proto := "http://"
	client := resty.New()

	_, err = client.R().
		Post(proto + url)

	if err != nil {
		return err
	}
	return nil
}
