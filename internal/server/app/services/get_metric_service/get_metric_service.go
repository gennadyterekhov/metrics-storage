package getmetricservice

import (
	"fmt"
	"github.com/gennadyterekhov/metrics-storage/internal/constants/exceptions"
	"github.com/gennadyterekhov/metrics-storage/internal/constants/types"
	"github.com/gennadyterekhov/metrics-storage/internal/domain/models"
	"github.com/gennadyterekhov/metrics-storage/internal/logger"
	"github.com/gennadyterekhov/metrics-storage/internal/server/storage"
	"strconv"
)

func GetMetricAsString(metricType string, name string) (metric string, err error) {
	logger.ZapSugarLogger.Debugln("GetMetricAsString name, metricType", name, metricType)

	if metricType == types.Counter {
		val, err := storage.MetricsRepository.GetCounter(name)
		if err != nil {
			return "", err
		}
		return strconv.FormatInt(val, 10), nil
	}
	if metricType == types.Gauge {
		val, err := storage.MetricsRepository.GetGauge(name)
		if err != nil {
			return "", err
		}
		return strconv.FormatFloat(val, 'g', -1, 64), nil
	}
	return "", fmt.Errorf(exceptions.InvalidMetricTypeChoice)
}

func GetMetricsAsStruct(metricType string, name string) (metric *models.Metrics, err error) {
	logger.ZapSugarLogger.Debugln("GetMetricsAsStruct name, metricType", name, metricType)

	metric = &models.Metrics{
		ID:    name,
		MType: metricType,
	}

	if metricType == types.Counter {
		val, err := storage.MetricsRepository.GetCounter(name)
		if err != nil {
			logger.ZapSugarLogger.Warnln("could not get counter by name", name, err.Error())
			return metric, err
		}
		metric.Delta = &val
		return metric, nil
	}
	if metricType == types.Gauge {
		val, err := storage.MetricsRepository.GetGauge(name)
		if err != nil {
			logger.ZapSugarLogger.Warnln("could not get gauge by name", name, err.Error())
			return metric, err
		}
		metric.Value = &val
		return metric, nil
	}
	return metric, fmt.Errorf(exceptions.InvalidMetricTypeChoice)
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
	for name, val := range storage.MetricsRepository.GetAllGauges() {
		list += fmt.Sprintf("<li>%v : %v</li>", name, val)
	}

	return list
}

func getCounterList() string {
	list := ""
	for name, val := range storage.MetricsRepository.GetAllCounters() {
		list += fmt.Sprintf("<li>%v : %v</li>", name, val)
	}
	return list
}
