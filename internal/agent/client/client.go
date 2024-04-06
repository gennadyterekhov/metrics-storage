package client

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"github.com/Rican7/retry"
	"github.com/Rican7/retry/backoff"
	"github.com/Rican7/retry/strategy"
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
	"time"
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

	metricToEncode := requests.SaveMetricListRequest{
		getSubrequest(&memStats.Alloc),
		getSubrequest(&memStats.BuckHashSys),
		getSubrequest(&memStats.Frees),
		getSubrequest(&memStats.GCCPUFraction),
		getSubrequest(&memStats.GCSys),
		getSubrequest(&memStats.HeapAlloc),
		getSubrequest(&memStats.HeapIdle),
		getSubrequest(&memStats.HeapInuse),
		getSubrequest(&memStats.HeapObjects),
		getSubrequest(&memStats.HeapReleased),
		getSubrequest(&memStats.HeapSys),
		getSubrequest(&memStats.LastGC),
		getSubrequest(&memStats.Lookups),
		getSubrequest(&memStats.MCacheInuse),
		getSubrequest(&memStats.MCacheSys),
		getSubrequest(&memStats.MSpanInuse),
		getSubrequest(&memStats.MSpanSys),
		getSubrequest(&memStats.Mallocs),
		getSubrequest(&memStats.NextGC),
		getSubrequest(&memStats.NumForcedGC),
		getSubrequest(&memStats.NumGC),
		getSubrequest(&memStats.OtherSys),
		getSubrequest(&memStats.PauseTotalNs),
		getSubrequest(&memStats.StackInuse),
		getSubrequest(&memStats.StackSys),
		getSubrequest(&memStats.Sys),
		getSubrequest(&memStats.TotalAlloc),
		getSubrequest(&memStats.PollCount),
		getSubrequest(&memStats.RandomValue),
	}
	jsonBytes, err := json.Marshal(metricToEncode)
	if err != nil {
		logger.ZapSugarLogger.Errorln("error when encoding metric batch", err.Error())

		return nil, err
	}
	return jsonBytes, nil
}

func getSubrequest(met metric.MetricURLFormatter) *requests.SaveMetricRequest {
	counter, gauge, err := getMetricValues(met)
	if err != nil {
		return nil
	}

	return &requests.SaveMetricRequest{
		MetricName:   met.GetName(),
		MetricType:   met.GetType(),
		CounterValue: &counter,
		GaugeValue:   &gauge,
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

	fullURL := domain
	if isBatch {
		fullURL += "/updates/"
	} else {
		fullURL += "/update/"
	}
	return fullURL
}

func sendBody(url string, body []byte) (err error) {
	request := client.R().
		SetBody(body).
		SetHeader(constants.HeaderContentType, constants.ApplicationJSON)

	err = sendRequestWithRetries(request, url)

	if err != nil {
		logger.ZapSugarLogger.Errorln("error when sending metric", err.Error())
		return err
	}

	return err
}

func sendBodyGzipCompressed(url string, body []byte) (err error) {
	request, err := prepareRequest(body)
	if err != nil {
		return err
	}
	err = sendRequestWithRetries(request, url)

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

func sendRequestWithRetries(request *resty.Request, url string) (err error) {

	err = retry.Retry(
		func(numberOfAttempt uint) error {
			_, err := request.Post(url)
			if err != nil {
				logger.ZapSugarLogger.Errorf(
					"error when sending request. attempt: %v error: %v",
					numberOfAttempt,
					err.Error(),
				)
				return err
			}
			return nil
		},
		strategy.Limit(3),
		strategy.Backoff(backoff.Incremental(0*time.Second, 3*time.Second)),
	)

	if err != nil {
		logger.ZapSugarLogger.Errorln("error when sending request with 3 retries", err.Error())
		return err
	}
	return nil
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
