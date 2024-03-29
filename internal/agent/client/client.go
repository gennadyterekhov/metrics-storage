package client

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"github.com/gennadyterekhov/metrics-storage/internal/agent/metric"
	"github.com/gennadyterekhov/metrics-storage/internal/constants"
	"github.com/gennadyterekhov/metrics-storage/internal/constants/types"
	"github.com/gennadyterekhov/metrics-storage/internal/domain/models"
	"github.com/gennadyterekhov/metrics-storage/internal/logger"
	"github.com/gennadyterekhov/metrics-storage/internal/server/httpui/requests"
	"github.com/go-resty/resty/v2"
	"net/http"
	"strconv"
	"strings"
)

var client *resty.Client

type MetricsStorageClient struct {
	Address string
	IsGzip  bool
}

func init() {
	client = resty.New()
}

func (msc *MetricsStorageClient) SendMetric(met metric.MetricURLFormatter) (err error) {
	jsonBytes, err := getBody(met)
	if err != nil {
		return err
	}

	err = msc.sendRequestToMetricsServer(jsonBytes, false)
	if err != nil {
		logger.ZapSugarLogger.Errorln("error when sending metric "+met.GetName()+" to server", err.Error())
		return err
	}
	return nil
}

func (msc *MetricsStorageClient) SendAllMetricsInOneRequest(memStats *metric.MetricsSet) (err error) {
	jsonBytes, err := getBodyForAllMetrics(memStats)
	if err != nil {
		return err
	}
	err = msc.sendRequestToMetricsServer(jsonBytes, true)

	if err != nil {
		logger.ZapSugarLogger.Errorln("error when sending metric batch to server", err.Error())
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

func getBodyForAllMetrics(memStats *metric.MetricsSet) ([]byte, error) {

	metricToEncode := requests.SaveMetricBatchRequest{
		Alloc:         getGaugeSubrequest(&memStats.Alloc),
		BuckHashSys:   getGaugeSubrequest(&memStats.BuckHashSys),
		Frees:         getGaugeSubrequest(&memStats.Frees),
		GCCPUFraction: getGaugeSubrequest(&memStats.GCCPUFraction),
		GCSys:         getGaugeSubrequest(&memStats.GCSys),
		HeapAlloc:     getGaugeSubrequest(&memStats.HeapAlloc),
		HeapIdle:      getGaugeSubrequest(&memStats.HeapIdle),
		HeapInuse:     getGaugeSubrequest(&memStats.HeapInuse),
		HeapObjects:   getGaugeSubrequest(&memStats.HeapObjects),
		HeapReleased:  getGaugeSubrequest(&memStats.HeapReleased),
		HeapSys:       getGaugeSubrequest(&memStats.HeapSys),
		LastGC:        getGaugeSubrequest(&memStats.LastGC),
		Lookups:       getGaugeSubrequest(&memStats.Lookups),
		MCacheInuse:   getGaugeSubrequest(&memStats.MCacheInuse),
		MCacheSys:     getGaugeSubrequest(&memStats.MCacheSys),
		MSpanInuse:    getGaugeSubrequest(&memStats.MSpanInuse),
		MSpanSys:      getGaugeSubrequest(&memStats.MSpanSys),
		Mallocs:       getGaugeSubrequest(&memStats.Mallocs),
		NextGC:        getGaugeSubrequest(&memStats.NextGC),
		NumForcedGC:   getGaugeSubrequest(&memStats.NumForcedGC),
		NumGC:         getGaugeSubrequest(&memStats.NumGC),
		OtherSys:      getGaugeSubrequest(&memStats.OtherSys),
		PauseTotalNs:  getGaugeSubrequest(&memStats.PauseTotalNs),
		StackInuse:    getGaugeSubrequest(&memStats.StackInuse),
		StackSys:      getGaugeSubrequest(&memStats.StackSys),
		Sys:           getGaugeSubrequest(&memStats.Sys),
		TotalAlloc:    getGaugeSubrequest(&memStats.TotalAlloc),
		PollCount:     getCounterSubrequest(&memStats.PollCount),
		RandomValue:   getGaugeSubrequest(&memStats.RandomValue),
	}
	jsonBytes, err := json.Marshal(metricToEncode)
	if err != nil {
		logger.ZapSugarLogger.Errorln("error when encoding metric batch", err.Error())

		return nil, err
	}
	return jsonBytes, nil
}

func getGaugeSubrequest(met metric.MetricURLFormatter) *requests.GaugeMetricSubrequest {
	_, gauge, err := getMetricValues(met)
	if err != nil {
		return nil
	}

	return &requests.GaugeMetricSubrequest{MetricName: met.GetName(), MetricType: met.GetType(), GaugeValue: gauge}
}

func getCounterSubrequest(met metric.MetricURLFormatter) *requests.CounterMetricSubrequest {
	counter, _, err := getMetricValues(met)
	if err != nil {
		return nil
	}

	return &requests.CounterMetricSubrequest{
		MetricName:   met.GetName(),
		MetricType:   met.GetType(),
		CounterValue: counter,
	}
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

func (msc *MetricsStorageClient) sendRequestToMetricsServer(body []byte, isBatch bool) (err error) {
	fullURL := getFullURL(msc.Address, isBatch)

	if msc.IsGzip {
		logger.ZapSugarLogger.Debugln("sending GZIP metric to server", http.MethodPost, fullURL, string(body))

		err = sendBodyGzipCompressed(fullURL, body)
	} else {
		logger.ZapSugarLogger.Debugln("sending metric to server", http.MethodPost, fullURL, string(body))

		err = sendBody(fullURL, body)
	}
	if err != nil {
		return err
	}
	return nil
}

func getFullURL(domain string, isBatch bool) string {
	proto := "http://"

	if !strings.Contains(domain, proto) {
		domain = proto + domain
	}

	fullURL := domain + "/update/"
	if isBatch {
		fullURL += "batch/"
	}
	return fullURL
}

func sendBody(url string, body []byte) (err error) {
	request := client.R().
		SetBody(body).
		SetHeader(constants.HeaderContentType, constants.ApplicationJSON)

	response, err := request.Post(url)
	if err != nil {
		logger.ZapSugarLogger.Errorln("error when sending metric", err.Error())
		return err
	}
	logger.ZapSugarLogger.Infoln("sending metric response", response)

	return err
}

func sendBodyGzipCompressed(url string, body []byte) (err error) {
	request, err := prepareRequest(body)
	if err != nil {
		return err
	}
	response, err := request.Post(url)
	logger.ZapSugarLogger.Infoln("server response", response)

	if err != nil {
		logger.ZapSugarLogger.Errorln("error when sending compressed metric", err.Error())
		return err
	}

	return err
}

func prepareRequest(body []byte) (*resty.Request, error) {
	compressedBody, err := getCompressedBody(body)
	if err != nil {
		return nil, err
	}
	request := client.R().
		SetHeader(constants.HeaderContentType, constants.ApplicationJSON).
		SetHeader("Accept-Encoding", "gzip").
		SetHeader("Content-Encoding", "gzip").
		SetBody(compressedBody)

	return request, nil
}

func getCompressedBody(body []byte) (*bytes.Buffer, error) {
	logger.ZapSugarLogger.Debugln(" body before compression as sent by agent", string(body))

	var bodyBuffer bytes.Buffer
	compressedBodyWriter, err := gzip.NewWriterLevel(&bodyBuffer, gzip.BestSpeed)
	if err != nil {
		logger.ZapSugarLogger.Errorln("error when opening gzip writer", err.Error())
		return nil, err
	}
	defer compressedBodyWriter.Close()
	_, err = compressedBodyWriter.Write(body)
	if err != nil {
		logger.ZapSugarLogger.Errorln("error when writing gzip body", err.Error())
		return nil, err
	}
	err = compressedBodyWriter.Flush()
	if err != nil {
		logger.ZapSugarLogger.Errorln("error when flushing gzip body", err.Error())
		return nil, err
	}
	logger.ZapSugarLogger.Debugln("compressed body as sent by agent", bodyBuffer.String())

	return &bodyBuffer, nil
}
