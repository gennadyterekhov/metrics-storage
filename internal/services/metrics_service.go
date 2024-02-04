package services

import (
	"fmt"
	"github.com/gennadyterekhov/metrics-storage/internal/container"
	"github.com/gennadyterekhov/metrics-storage/internal/types"
	"strconv"
)

func SaveMetricToMemory(metricType string, name string, counterValue int64, gaugeValue float64) {
	if metricType == types.Counter {
		container.Instance.MetricsRepository.AddCounter(name, counterValue)
	}
	if metricType == types.Gauge {
		container.Instance.MetricsRepository.AddGauge(name, gaugeValue)
	}
}

func GetMetricAsString(metricType string, name string) string {
	if metricType == types.Counter {
		val := container.Instance.MetricsRepository.GetCounter(name)

		return strconv.FormatInt(val, 10)
	}
	if metricType == types.Gauge {
		val := container.Instance.MetricsRepository.GetGauge(name)

		return strconv.FormatFloat(val, 'E', 2, 64)
	}
	return ""
}

func GetMetricsListAsHTML() string {
	templateText := `
<!DOCTYPE html>
<html>
  <head></head>
  <body>
    <h2>%v</h2>
    <ul>
%v
    </ul>
    <h2>%v</h2>
    <ul>
%v
    </ul>
  </body>
</html>
`
	gaugeList := getGaugeList()
	counterList := getCounterList()
	return fmt.Sprintf(
		templateText,
		types.Gauge,
		gaugeList,
		types.Counter,
		counterList,
	)
}

func getGaugeList() string {
	list := ""
	for name, val := range container.Instance.MetricsRepository.GetAllGauges() {
		list += fmt.Sprintf("<li>%v : %v</li>", name, val)
	}

	return list
}

func getCounterList() string {
	list := ""
	for name, val := range container.Instance.MetricsRepository.GetAllCounters() {
		list += fmt.Sprintf("<li>%v : %v</li>", name, val)
	}
	return list
}
