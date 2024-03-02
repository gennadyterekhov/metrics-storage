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
	"github.com/go-resty/resty/v2"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var client *resty.Client

type MetricsStorageClient struct {
	Address          string
	IsGzip           bool
	SendWhenNoServer bool
}

func init() {
	client = resty.New()
}

func (msc *MetricsStorageClient) SendMetric(met metric.MetricURLFormatter) (err error) {
	msc.SendWhenNoServer = true
	jsonBytes, err := getBody(met)
	if err != nil {
		return err
	}

	err = msc.sendRequestToMetricsServer(jsonBytes)
	if err != nil {
		logger.ZapSugarLogger.Warnln("error when sending metric "+met.GetName()+" to server", err.Error())
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

func (msc *MetricsStorageClient) sendRequestToMetricsServer(body []byte) (err error) {
	fullURL := getFullURL(msc.Address)

	if msc.IsGzip {
		logger.ZapSugarLogger.Debugln("sending GZIP metric to server", http.MethodPost, fullURL, string(body))

		err = sendBodyGzipCompressed(fullURL, body, msc.SendWhenNoServer)
	} else {
		logger.ZapSugarLogger.Debugln("sending metric to server", http.MethodPost, fullURL, string(body))

		err = sendBody(fullURL, body, msc.SendWhenNoServer)
	}
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

func sendBody(url string, body []byte, sendWhenNoServer bool) (err error) {
	request := client.R().
		SetBody(body).
		SetHeader(constants.HeaderContentType, constants.ApplicationJSON)

	response, err := request.Post(url)
	for err != nil {
		logger.ZapSugarLogger.Warnln("error when sending metric", err.Error())

		if !sendWhenNoServer {
			break
		}
		time.Sleep(time.Millisecond * 50)
		_, err = request.Post(url)
	}
	logger.ZapSugarLogger.Infoln("sending metric response", response)

	return err
}

func sendBodyGzipCompressed(url string, body []byte, sendWhenNoServer bool) (err error) {
	request, err := prepareRequest(body)
	if err != nil {
		return err
	}
	response, err := request.Post(url)
	logger.ZapSugarLogger.Infoln("server response", response)

	for err != nil {
		logger.ZapSugarLogger.Warnln("error when sending compressed metric", err.Error())

		if !sendWhenNoServer {
			break
		}
		time.Sleep(time.Millisecond * 50)
		_, err = request.Post(url)
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
	defer func(compressedBodyWriter *gzip.Writer) {
		err := compressedBodyWriter.Close()
		if err != nil {
			logger.ZapSugarLogger.Debugln("could not close body", err.Error())
		}
	}(compressedBodyWriter)

	return bodyBuffer, nil
}
