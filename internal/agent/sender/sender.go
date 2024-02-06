package sender

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"runtime"
	"sync"
	"time"
)

type Sender interface {
	Send() error
	wait()
	shouldContinue(int)
}

type MetricsSender struct {
	Address     string
	Interval    int
	Channel     chan runtime.MemStats
	IsRunning   bool
	IsRunningMu *sync.Mutex
}

func (msnd *MetricsSender) shouldContinue(iter int) bool {
	//return iter == 0
	return true
}

func (msnd *MetricsSender) wait() {
	time.Sleep(time.Duration(msnd.Interval * int(time.Second)))
}

func (msnd *MetricsSender) reportRoutine(memStats *runtime.MemStats) {

	sendAllMetrics(msnd.Address, memStats, 1)
}

func (msnd *MetricsSender) Report(memStatsPtr *runtime.MemStats) {
	msnd.IsRunningMu.Lock()
	msnd.IsRunning = true
	fmt.Println("msnd.IsRunning", msnd.IsRunning)

	msnd.wait()

	fmt.Println("reporting runtime metrics, getting from channel")
	//memStats := <-msnd.Channel
	//fmt.Println("memStats", memStats)
	fmt.Println("GOT from channel")

	msnd.reportRoutine(memStatsPtr)
	msnd.IsRunning = false
	msnd.IsRunningMu.Unlock()
}

func sendAllMetrics(address string, memStats *runtime.MemStats, pollCount int) {
	fmt.Println("sending all metrics ")

	urls := getURLs(memStats)
	for i := 0; i < len(urls); i++ {
		_ = sendMetric(address + urls[i])
	}
	_ = sendMetric(address + fmt.Sprintf("/update/counter/PollCount/%v", pollCount))
	_ = sendMetric(address + fmt.Sprintf("/update/gauge/RandomValue/%v", pollCount))
}

func getURLs(memStats *runtime.MemStats) []string {
	return []string{
		getURL("Alloc", float64(memStats.Alloc)),
		getURL("BuckHashSys", float64(memStats.BuckHashSys)),
		getURL("Frees", float64(memStats.Frees)),
		getURL("GCCPUFraction", float64(memStats.GCCPUFraction)),
		getURL("GCSys", float64(memStats.GCSys)),
		getURL("HeapAlloc", float64(memStats.HeapAlloc)),
		getURL("HeapIdle", float64(memStats.HeapIdle)),
		getURL("HeapInuse", float64(memStats.HeapInuse)),
		getURL("HeapObjects", float64(memStats.HeapObjects)),
		getURL("HeapReleased", float64(memStats.HeapReleased)),
		getURL("HeapSys", float64(memStats.HeapSys)),
		getURL("LastGC", float64(memStats.LastGC)),
		getURL("Lookups", float64(memStats.Lookups)),
		getURL("MCacheInuse", float64(memStats.MCacheInuse)),
		getURL("MCacheSys", float64(memStats.MCacheSys)),
		getURL("MSpanInuse", float64(memStats.MSpanInuse)),
		getURL("MSpanSys", float64(memStats.MSpanSys)),
		getURL("Mallocs", float64(memStats.Mallocs)),
		getURL("NextGC", float64(memStats.NextGC)),
		getURL("NumForcedGC", float64(memStats.NumForcedGC)),
		getURL("NumGC", float64(memStats.NumGC)),
		getURL("OtherSys", float64(memStats.OtherSys)),
		getURL("PauseTotalNs", float64(memStats.PauseTotalNs)),
		getURL("StackInuse", float64(memStats.StackInuse)),
		getURL("StackSys", float64(memStats.StackSys)),
		getURL("Sys", float64(memStats.Sys)),
		getURL("TotalAlloc", float64(memStats.TotalAlloc)),
	}
}

func getURL(name string, val float64) string {
	template := "/update/gauge/%v/%v"
	return fmt.Sprintf(template, name, val)
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
