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
		logger.ZapSugarLogger.Warnln("error when sending metric "+met.GetName()+" to server", err.Error())
		return err
	}
	return nil
}

func (msc *MetricsStorageClient) SendAllMetricsInOneRequest(memStats *metric.MetricsSet) (err error) {
	jsonBytes, err := getBodyForAllMetrics(memStats)
	if err != nil {
		return err
	}

	err = retry.Retry(
		func(attempt uint) error {
			return msc.sendRequestToMetricsServer(jsonBytes, true)
		},
		strategy.Limit(4),
		strategy.Backoff(backoff.Incremental(-1*time.Second, 2*time.Second)),
	)

	if err != nil {
		logger.ZapSugarLogger.Warnln("error when sending metric batch to server", err.Error())
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
		logger.ZapSugarLogger.Warnln("error when encoding metric", err.Error())

		return nil, err
	}
	return jsonBytes, nil
}

func getBodyForAllMetrics(memStats *metric.MetricsSet) ([]byte, error) {

	metricToEncode := requests.SaveMetricBatchRequest{
		Alloc:         requests.GaugeMetricSubrequest{MetricName: memStats.Alloc.Name, MetricType: memStats.Alloc.Type, GaugeValue: memStats.Alloc.Value},
		BuckHashSys:   requests.GaugeMetricSubrequest{MetricName: memStats.BuckHashSys.Name, MetricType: memStats.BuckHashSys.Type, GaugeValue: memStats.BuckHashSys.Value},
		Frees:         requests.GaugeMetricSubrequest{MetricName: memStats.Frees.Name, MetricType: memStats.Frees.Type, GaugeValue: memStats.Frees.Value},
		GCCPUFraction: requests.GaugeMetricSubrequest{MetricName: memStats.GCCPUFraction.Name, MetricType: memStats.GCCPUFraction.Type, GaugeValue: memStats.GCCPUFraction.Value},
		GCSys:         requests.GaugeMetricSubrequest{MetricName: memStats.GCSys.Name, MetricType: memStats.GCSys.Type, GaugeValue: memStats.GCSys.Value},
		HeapAlloc:     requests.GaugeMetricSubrequest{MetricName: memStats.HeapAlloc.Name, MetricType: memStats.HeapAlloc.Type, GaugeValue: memStats.HeapAlloc.Value},
		HeapIdle:      requests.GaugeMetricSubrequest{MetricName: memStats.HeapIdle.Name, MetricType: memStats.HeapIdle.Type, GaugeValue: memStats.HeapIdle.Value},
		HeapInuse:     requests.GaugeMetricSubrequest{MetricName: memStats.HeapInuse.Name, MetricType: memStats.HeapInuse.Type, GaugeValue: memStats.HeapInuse.Value},
		HeapObjects:   requests.GaugeMetricSubrequest{MetricName: memStats.HeapObjects.Name, MetricType: memStats.HeapObjects.Type, GaugeValue: memStats.HeapObjects.Value},
		HeapReleased:  requests.GaugeMetricSubrequest{MetricName: memStats.HeapReleased.Name, MetricType: memStats.HeapReleased.Type, GaugeValue: memStats.HeapReleased.Value},
		HeapSys:       requests.GaugeMetricSubrequest{MetricName: memStats.HeapSys.Name, MetricType: memStats.HeapSys.Type, GaugeValue: memStats.HeapSys.Value},
		LastGC:        requests.GaugeMetricSubrequest{MetricName: memStats.LastGC.Name, MetricType: memStats.LastGC.Type, GaugeValue: memStats.LastGC.Value},
		Lookups:       requests.GaugeMetricSubrequest{MetricName: memStats.Lookups.Name, MetricType: memStats.Lookups.Type, GaugeValue: memStats.Lookups.Value},
		MCacheInuse:   requests.GaugeMetricSubrequest{MetricName: memStats.MCacheInuse.Name, MetricType: memStats.MCacheInuse.Type, GaugeValue: memStats.MCacheInuse.Value},
		MCacheSys:     requests.GaugeMetricSubrequest{MetricName: memStats.MCacheSys.Name, MetricType: memStats.MCacheSys.Type, GaugeValue: memStats.MCacheSys.Value},
		MSpanInuse:    requests.GaugeMetricSubrequest{MetricName: memStats.MSpanInuse.Name, MetricType: memStats.MSpanInuse.Type, GaugeValue: memStats.MSpanInuse.Value},
		MSpanSys:      requests.GaugeMetricSubrequest{MetricName: memStats.MSpanSys.Name, MetricType: memStats.MSpanSys.Type, GaugeValue: memStats.MSpanSys.Value},
		Mallocs:       requests.GaugeMetricSubrequest{MetricName: memStats.Mallocs.Name, MetricType: memStats.Mallocs.Type, GaugeValue: memStats.Mallocs.Value},
		NextGC:        requests.GaugeMetricSubrequest{MetricName: memStats.NextGC.Name, MetricType: memStats.NextGC.Type, GaugeValue: memStats.NextGC.Value},
		NumForcedGC:   requests.GaugeMetricSubrequest{MetricName: memStats.NumForcedGC.Name, MetricType: memStats.NumForcedGC.Type, GaugeValue: memStats.NumForcedGC.Value},
		NumGC:         requests.GaugeMetricSubrequest{MetricName: memStats.NumGC.Name, MetricType: memStats.NumGC.Type, GaugeValue: memStats.NumGC.Value},
		OtherSys:      requests.GaugeMetricSubrequest{MetricName: memStats.OtherSys.Name, MetricType: memStats.OtherSys.Type, GaugeValue: memStats.OtherSys.Value},
		PauseTotalNs:  requests.GaugeMetricSubrequest{MetricName: memStats.PauseTotalNs.Name, MetricType: memStats.PauseTotalNs.Type, GaugeValue: memStats.PauseTotalNs.Value},
		StackInuse:    requests.GaugeMetricSubrequest{MetricName: memStats.StackInuse.Name, MetricType: memStats.StackInuse.Type, GaugeValue: memStats.StackInuse.Value},
		StackSys:      requests.GaugeMetricSubrequest{MetricName: memStats.StackSys.Name, MetricType: memStats.StackSys.Type, GaugeValue: memStats.StackSys.Value},
		Sys:           requests.GaugeMetricSubrequest{MetricName: memStats.Sys.Name, MetricType: memStats.Sys.Type, GaugeValue: memStats.Sys.Value},
		TotalAlloc:    requests.GaugeMetricSubrequest{MetricName: memStats.TotalAlloc.Name, MetricType: memStats.TotalAlloc.Type, GaugeValue: memStats.TotalAlloc.Value},
		PollCount:     requests.CounterMetricSubrequest{MetricName: memStats.PollCount.Name, MetricType: memStats.PollCount.Type, CounterValue: memStats.PollCount.Value},
		RandomValue:   requests.GaugeMetricSubrequest{MetricName: memStats.RandomValue.Name, MetricType: memStats.RandomValue.Type, GaugeValue: memStats.RandomValue.Value},
	}
	jsonBytes, err := json.Marshal(metricToEncode)
	if err != nil {
		logger.ZapSugarLogger.Warnln("error when encoding metric batch", err.Error())

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
		logger.ZapSugarLogger.Warnln("error when sending metric", err.Error())
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
		logger.ZapSugarLogger.Warnln("error when sending compressed metric", err.Error())
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
	bodyBuffer := bytes.NewBuffer(body)
	compressedBodyWriter, err := gzip.NewWriterLevel(bodyBuffer, gzip.BestSpeed)
	if err != nil {
		logger.ZapSugarLogger.Warnln("error when opening gzip writer", err.Error())
		return nil, err
	}
	defer compressedBodyWriter.Close()
	_, err = compressedBodyWriter.Write(body)
	if err != nil {
		logger.ZapSugarLogger.Warnln("error when writing gzip body", err.Error())
		return nil, err
	}
	err = compressedBodyWriter.Flush()
	if err != nil {
		logger.ZapSugarLogger.Warnln("error when flushing gzip body", err.Error())
		return nil, err
	}

	return bodyBuffer, nil
}
