package agent

import (
	"fmt"
	"net/http"
	"runtime"
	"time"
)

func Agent() (err error) {
	//pollInterval := 2
	reportInterval := 10
	memStats := &runtime.MemStats{}

	address := `http://localhost:8080`
	var urls []string
	for i := 0; ; i++ {

		runtime.ReadMemStats(memStats)

		urls = getUrls(memStats)

		for i := 0; i < len(urls); i++ {
			err = sendMetric(address + urls[i])
			if err != nil {
				return err
			}
		}
		err = sendMetric(address + fmt.Sprintf("/update/counter/PollCount/%v", i))
		if err != nil {
			return err
		}
		err = sendMetric(address + fmt.Sprintf("/update/gauge/RandomValue/%v", i))
		if err != nil {
			return err
		}

		//time.Sleep(time.Duration(pollInterval * int(time.Second)))
		time.Sleep(time.Duration(reportInterval * int(time.Second)))
	}

	return nil
}

func getUrls(memStats *runtime.MemStats) []string {
	return []string{
		getUrl("Alloc", float64(memStats.Alloc)),
		getUrl("BuckHashSys", float64(memStats.BuckHashSys)),
		getUrl("Frees", float64(memStats.Frees)),
		getUrl("GCCPUFraction", float64(memStats.GCCPUFraction)),
		getUrl("GCSys", float64(memStats.GCSys)),
		getUrl("HeapAlloc", float64(memStats.HeapAlloc)),
		getUrl("HeapIdle", float64(memStats.HeapIdle)),
		getUrl("HeapInuse", float64(memStats.HeapInuse)),
		getUrl("HeapObjects", float64(memStats.HeapObjects)),
		getUrl("HeapReleased", float64(memStats.HeapReleased)),
		getUrl("HeapSys", float64(memStats.HeapSys)),
		getUrl("LastGC", float64(memStats.LastGC)),
		getUrl("Lookups", float64(memStats.Lookups)),
		getUrl("MCacheInuse", float64(memStats.MCacheInuse)),
		getUrl("MCacheSys", float64(memStats.MCacheSys)),
		getUrl("MSpanInuse", float64(memStats.MSpanInuse)),
		getUrl("MSpanSys", float64(memStats.MSpanSys)),
		getUrl("Mallocs", float64(memStats.Mallocs)),
		getUrl("NextGC", float64(memStats.NextGC)),
		getUrl("NumForcedGC", float64(memStats.NumForcedGC)),
		getUrl("NumGC", float64(memStats.NumGC)),
		getUrl("OtherSys", float64(memStats.OtherSys)),
		getUrl("PauseTotalNs", float64(memStats.PauseTotalNs)),
		getUrl("StackInuse", float64(memStats.StackInuse)),
		getUrl("StackSys", float64(memStats.StackSys)),
		getUrl("Sys", float64(memStats.Sys)),
		getUrl("TotalAlloc", float64(memStats.TotalAlloc)),
	}
}

func getUrl(name string, val float64) string {
	template := "/update/gauge/%v/%v"
	return fmt.Sprintf(template, name, val)
}

func sendMetric(url string) (err error) {
	_, err = http.Post(url, "text/plain", nil)
	if err != nil {
		return err
	}
	return nil
}
