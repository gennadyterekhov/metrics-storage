package bodypreparer

import (
	"bytes"
	"compress/gzip"
	"encoding/hex"
	"github.com/gennadyterekhov/metrics-storage/internal/constants"
	"github.com/gennadyterekhov/metrics-storage/internal/helper/hasher"
	"github.com/gennadyterekhov/metrics-storage/internal/logger"
	"github.com/go-resty/resty/v2"
)

func PrepareRequest(client *resty.Client, body []byte, isGzip bool, key string) (*resty.Request, error) {
	request := client.R().
		SetHeader(constants.HeaderContentType, constants.ApplicationJSON)

	if key != "" {
		request.SetHeader("HashSHA256", hex.EncodeToString(hashBytes(body, key)))
	}

	if !isGzip {
		request.SetBody(body)
		return request, nil
	}

	request, err := prepareCompressed(request, body)
	if err != nil {
		return nil, err
	}
	return request, nil
}

func prepareCompressed(request *resty.Request, body []byte) (*resty.Request, error) {
	compressedBody, err := getCompressedBody(body)
	if err != nil {
		return nil, err
	}
	request.SetHeader("Accept-Encoding", "gzip").
		SetHeader("Content-Encoding", "gzip").
		SetBody(compressedBody)

	return request, nil
}

func hashBytes(target []byte, key string) []byte {
	return hasher.HashBytes(target, key)
}

func getCompressedBody(body []byte) (*bytes.Buffer, error) {
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
