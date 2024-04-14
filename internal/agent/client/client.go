package client

import (
	"github.com/Rican7/retry"
	"github.com/Rican7/retry/backoff"
	"github.com/Rican7/retry/strategy"
	"github.com/gennadyterekhov/metrics-storage/internal/agent/client/bodymaker"
	"github.com/gennadyterekhov/metrics-storage/internal/agent/client/bodypreparer"
	"github.com/gennadyterekhov/metrics-storage/internal/agent/metric"
	"github.com/gennadyterekhov/metrics-storage/internal/logger"
	"github.com/go-resty/resty/v2"
	"net/http"
	"strings"
	"time"
)

var client *resty.Client

type MetricsStorageClient struct {
	Address             string
	IsGzip              bool
	PayloadSignatureKey string
}

func init() {
	client = resty.New()
}

func (msc *MetricsStorageClient) SendMetric(met metric.MetricURLFormatter) (err error) {
	jsonBytes, err := bodymaker.GetBody(met)
	if err != nil {
		return err
	}

	err = msc.sendRequestToMetricsServer(jsonBytes, false)
	if err != nil {
		logger.ZapSugarLogger.Errorln("error when sending metric "+met.GetName()+" to server", err.Error())
		return err
	}
	logger.ZapSugarLogger.Debugln("request seemingly sent without errors")

	return nil
}

func (msc *MetricsStorageClient) SendAllMetricsInOneRequest(memStats *metric.MetricsSet) (err error) {
	jsonBytes, err := bodymaker.GetBodyForAllMetrics(memStats)
	if err != nil {
		return err
	}
	err = msc.sendRequestToMetricsServer(jsonBytes, true)

	if err != nil {
		logger.ZapSugarLogger.Errorln("error when sending metric batch to server", err.Error())
		return err
	}
	logger.ZapSugarLogger.Debugln("request seemingly sent without errors")

	return nil
}

func (msc *MetricsStorageClient) sendRequestToMetricsServer(body []byte, isBatch bool) (err error) {
	fullURL := getFullURL(msc.Address, isBatch)

	if msc.IsGzip {
		logger.ZapSugarLogger.Debugln("sending GZIP metric to server", http.MethodPost, fullURL, string(body))

		err = sendBodyGzipCompressed(fullURL, body, msc.PayloadSignatureKey)
	} else {
		logger.ZapSugarLogger.Debugln("sending metric to server", http.MethodPost, fullURL, string(body))

		err = sendBody(fullURL, body, msc.PayloadSignatureKey)
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

func sendBody(url string, body []byte, key string) (err error) {
	request, err := bodypreparer.PrepareRequest(client, body, false, key)
	if err != nil {
		logger.ZapSugarLogger.Errorln("error when preparing request", err.Error())
		return err
	}

	err = sendRequestWithRetries(request, url)

	if err != nil {
		logger.ZapSugarLogger.Errorln("error when sending metric", err.Error())
		return err
	}

	return err
}

func sendBodyGzipCompressed(url string, body []byte, key string) (err error) {
	request, err := bodypreparer.PrepareRequest(client, body, true, key)
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

func sendRequestWithRetries(request *resty.Request, url string) (err error) {

	err = retry.Retry(
		func(numberOfAttempt uint) error {
			_, err := request.Post(url)
			if err != nil {
				logger.ZapSugarLogger.Debugf(
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
