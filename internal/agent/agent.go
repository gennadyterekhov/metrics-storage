package agent

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"runtime"
	"time"
)

type shouldContinueType func(int) bool

func Agent(address string, shouldContinue shouldContinueType) (err error) {
	fmt.Println("Agent")

	//pollInterval := 2
	reportInterval := 10
	memStats := &runtime.MemStats{}

	var urls []string
	for i := 0; shouldContinue(i); i++ {

		runtime.ReadMemStats(memStats)

		urls = getURLs(memStats)

		for i := 0; i < len(urls); i++ {
			fmt.Println("before sendMetric")

			err = sendMetric(address + urls[i])
			if err != nil {
				fmt.Println("error from sendMetric, ignoring it")

				//return err
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

		if shouldContinue(i + 1) {
			//time.Sleep(time.Duration(pollInterval * int(time.Second)))
			time.Sleep(time.Duration(reportInterval * int(time.Second)))
		}
	}
	return nil
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
	fmt.Println("func sendMetric")

	client := resty.New()

	resp, err := client.R().
		Post(url)

	fmt.Println(resp.Body())

	fmt.Println(err)

	if err != nil {
		return err
	}
	return nil
}
