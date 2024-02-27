package client

import (
	"encoding/json"
	"github.com/gennadyterekhov/metrics-storage/internal/agent/metric"
	"github.com/gennadyterekhov/metrics-storage/internal/constants"
	"github.com/gennadyterekhov/metrics-storage/internal/constants/types"
	"github.com/gennadyterekhov/metrics-storage/internal/domain/models"
	"github.com/gennadyterekhov/metrics-storage/internal/logger"
	"github.com/go-resty/resty/v2"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func SendMetric(met metric.MetricURLFormatter, domain string) (err error) {
	jsonBytes, err := getBody(met)
	if err != nil {
		return err
	}

	err = sendRequestToMetricsServer(domain, jsonBytes)
	if err != nil {
		logger.ZapSugarLogger.Errorln("error when sending metrics to server", err.Error())
		return err
	}
	return nil
}

func getBody(met metric.MetricURLFormatter) ([]byte, error) {
	counterValue, gaugeValue, err := getMetricValues(met)
	if err != nil {
		return nil, err
	}

	metricToEncode := models.Metrics{
		ID:    met.GetName(),
		MType: met.GetType(),
		Delta: &counterValue,
		Value: &gaugeValue,
	}
	jsonBytes, err := json.Marshal(metricToEncode)
	if err != nil {
		logger.ZapSugarLogger.Errorln("error when encoding metric", err.Error())

		return nil, err
	}
	return jsonBytes, nil
}

func getMetricValues(met metric.MetricURLFormatter) (counterValue int64, gaugeValue float64, err error) {
	if met.GetType() == types.Counter {
		counterValue, err = strconv.ParseInt(met.GetValueAsString(), 10, 64)
		if err != nil {
			return 0, 0, err
		}
	} else {
		gaugeValue, err = strconv.ParseFloat(met.GetValueAsString(), 64)
		if err != nil {
			return 0, 0, err
		}
	}

	return counterValue, gaugeValue, nil
}

func sendRequestToMetricsServer(domain string, body []byte) (err error) {
	fullURL := getFullURL(domain)
	logger.ZapSugarLogger.Debugln("sending metric to server", http.MethodPost, fullURL, constants.ApplicationJSON, string(body))

	err = sendBody(fullURL, body)
	if err != nil {
		return err
	}
	return nil
}

func getFullURL(domain string) string {
	proto := "http://"

	if !strings.Contains(domain, proto) {
		domain = proto + domain
	}

	fullURL := domain + "/update/"
	return fullURL
}

func sendBody(url string, body []byte) (err error) {

	client := resty.New()
	_, err = client.R().
		SetBody(body).
		SetHeader(constants.HeaderContentType, constants.ApplicationJSON).
		Post(url)
	for err != nil {
		time.Sleep(time.Second)
		_, err = client.R().
			SetBody(body).
			SetHeader(constants.HeaderContentType, constants.ApplicationJSON).
			Post(url)
	}

	return err
}
