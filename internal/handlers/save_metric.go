package handlers

import (
	"errors"
	"fmt"
	"github.com/gennadyterekhov/metrics-storage/internal/container"
	"github.com/gennadyterekhov/metrics-storage/internal/types"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func SaveMetric(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Method not allowed. Allowed: POST", http.StatusMethodNotAllowed)
		return
	}

	var err error = nil
	metricTypeRaw, nameRaw, valueRaw, err := parseURL(req.URL)
	if err != nil && err.Error() == "expected exactly 3 parameters" {
		http.Error(res, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	counterValue, gaugeValue, err := validateParameters(metricTypeRaw, nameRaw, valueRaw)
	if err != nil && err.Error() == "name must be a non empty string" {
		http.Error(res, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if metricTypeRaw == types.Counter {
		container.Instance.MemStorage.AddCounter(nameRaw, counterValue)
	}
	if metricTypeRaw == types.Gauge {
		container.Instance.MemStorage.AddGauge(nameRaw, gaugeValue)
	}
}

func parseURL(url *url.URL) (metricType string, name string, value string, err error) {
	parameters := strings.Split(url.Path, "/")

	if len(parameters) != 5 { // empty string before first slash, update and 3 params
		return "", "", "", errors.New("expected exactly 3 parameters")
	}
	return parameters[2], parameters[3], parameters[4], nil
}

func validateParameters(metricTypeRaw string, nameRaw string, valueRaw string) (counterValue int64, gaugeValue float64, err error) {
	err = validateMetricType(metricTypeRaw)
	if err != nil {
		return 0, 0, err
	}
	err = validateMetricName(nameRaw)
	if err != nil {
		return 0, 0, err
	}
	counterValue, gaugeValue, err = validateMetricValue(metricTypeRaw, valueRaw)
	if err != nil {
		return 0, 0, err
	}

	return counterValue, gaugeValue, nil
}

func validateMetricType(metricTypeRaw string) (err error) {
	if metricTypeRaw != types.Counter && metricTypeRaw != types.Gauge {
		return fmt.Errorf("type can be %v or %v", types.Counter, types.Gauge)
	}
	return nil
}

func validateMetricName(nameRaw string) (err error) {
	if len(nameRaw) < 1 {
		return errors.New("name must be a non empty string")
	}
	return nil
}

func validateMetricValue(metricTypeValidated string, valueRaw string) (counterValue int64, gaugeValue float64, err error) {
	if metricTypeValidated == types.Counter {
		val, err := strconv.ParseInt(valueRaw, 10, 64)
		if err != nil {
			return 0, 0, nil
		}
		return val, 0, nil
	}
	if metricTypeValidated == types.Gauge {
		val, err := strconv.ParseFloat(valueRaw, 64)
		if err != nil {
			return 0, 0, nil
		}
		return 0, val, nil
	}
	return 0, 0, errors.New("unexpected type after validation")
}
