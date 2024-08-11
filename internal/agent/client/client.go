package client

import (
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/Rican7/retry"
	"github.com/Rican7/retry/backoff"
	"github.com/Rican7/retry/strategy"
	"github.com/gennadyterekhov/metrics-storage/internal/agent/client/bodymaker"
	"github.com/gennadyterekhov/metrics-storage/internal/agent/client/bodypreparer"
	"github.com/gennadyterekhov/metrics-storage/internal/agent/metric"
	"github.com/gennadyterekhov/metrics-storage/internal/common/logger"
	"github.com/go-resty/resty/v2"
)

type MetricsStorageClient struct {
	Address             string
	IsGzip              bool
	PayloadSignatureKey string
	RestyClient         *resty.Client
	PublicKeyFilePath   string
}

func (msc *MetricsStorageClient) SendMetric(met metric.URLFormatter) (err error) {
	jsonBytes, err := bodymaker.GetBody(met, msc.PublicKeyFilePath)
	if err != nil {
		return err
	}

	err = msc.sendRequestToMetricsServer(jsonBytes, false)
	if err != nil {
		return errors.Wrap(err, "error when sending metric "+met.GetName()+" to server")
	}

	return nil
}

func (msc *MetricsStorageClient) SendAllMetricsInOneRequest(memStats *metric.MetricsSet) (err error) {
	jsonBytes, err := bodymaker.GetBodyForAllMetrics(memStats)
	if err != nil {
		return err
	}
	err = msc.sendRequestToMetricsServer(jsonBytes, true)
	if err != nil {
		return errors.Wrap(err, "error when sending metric batch to server")
	}
	logger.Custom.Debugln("request seemingly sent without errors")

	return nil
}

func (msc *MetricsStorageClient) sendRequestToMetricsServer(body []byte, isBatch bool) (err error) {
	fullURL := getFullURL(msc.Address, isBatch)

	if msc.IsGzip {
		logger.Custom.Debugln("sending GZIP metric to server", http.MethodPost, fullURL, string(body))

		err = sendBodyGzipCompressed(msc.RestyClient, fullURL, body, msc.PayloadSignatureKey)
	} else {
		logger.Custom.Debugln("sending metric to server", http.MethodPost, fullURL, string(body))

		err = sendBody(msc.RestyClient, fullURL, body, msc.PayloadSignatureKey)
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

func sendBody(client *resty.Client, url string, body []byte, key string) (err error) {
	request, err := bodypreparer.PrepareRequest(client, body, false, key)
	if err != nil {
		return errors.Wrap(err, "error when preparing request")
	}

	err = sendRequestWithRetries(request, url)
	if err != nil {
		return errors.Wrap(err, "error when sending metric")
	}

	return err
}

func sendBodyGzipCompressed(client *resty.Client, url string, body []byte, key string) (err error) {
	request, err := bodypreparer.PrepareRequest(client, body, true, key)
	if err != nil {
		return err
	}
	err = sendRequestWithRetries(request, url)
	if err != nil {
		return errors.Wrap(err, "error when sending compressed metric")
	}

	return err
}

func sendRequestWithRetries(request *resty.Request, url string) (err error) {
	err = retry.Retry(
		attempt(request, url),
		strategy.Limit(3),
		strategy.Backoff(backoff.Incremental(0*time.Second, 3*time.Second)),
	)
	if err != nil {
		return errors.Wrap(err, "error when sending request with 3 retries")
	}
	return nil
}

func attempt(request *resty.Request, url string) func(numberOfAttempt uint) error {
	return func(numberOfAttempt uint) error {
		_, err := request.Post(url)
		if err != nil {
			return errors.Wrapf(err, "error when sending request. attempt: %v", numberOfAttempt)
		}
		return nil
	}
}
