package services

import (
	"fmt"
	"github.com/gennadyterekhov/metrics-storage/internal/container"
	"github.com/gennadyterekhov/metrics-storage/internal/exceptions"
	"github.com/gennadyterekhov/metrics-storage/internal/types"
	"strconv"
)

func SaveMetricToMemory(metricType string, name string, counterValue int64, gaugeValue float64) {
	fmt.Println("saving metrics " + name)
	if metricType == types.Counter {
		container.MetricsRepository.AddCounter(name, counterValue)
	}
	if metricType == types.Gauge {
		container.MetricsRepository.SetGauge(name, gaugeValue)
	}
}

func GetMetricAsString(metricType string, name string) (metric string, err error) {
	if metricType == types.Counter {
		val, err := container.MetricsRepository.GetCounter(name)
		if err != nil {
			return "", err
		}
		return strconv.FormatInt(val, 10), nil
	}
	if metricType == types.Gauge {
		val, err := container.MetricsRepository.GetGauge(name)
		if err != nil {
			return "", err
		}
		return strconv.FormatFloat(val, 'g', -1, 64), nil
	}
	return "", fmt.Errorf(exceptions.InvalidMetricTypeChoice)
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
	for name, val := range container.MetricsRepository.GetAllGauges() {
		list += fmt.Sprintf("<li>%v : %v</li>", name, val)
	}

	return list
}

func getCounterList() string {
	list := ""
	for name, val := range container.MetricsRepository.GetAllCounters() {
		list += fmt.Sprintf("<li>%v : %v</li>", name, val)
	}
	return list
}
