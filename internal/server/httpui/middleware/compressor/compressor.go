package compressor

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/gennadyterekhov/metrics-storage/internal/common/logger"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func GzipCompressor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if isRequestCompressed(request) {
			request = decompressRequest(request)
		}

		if isGzipAvailableForThisRequest(request) {
			compressionWriter := getCompressedWriter(response)
			if compressionWriter == nil {
				return
			}
			defer func(compressionWriter *gzip.Writer) {
				err := compressionWriter.Flush()
				if err != nil {
					logger.ZapSugarLogger.Errorln("error when flushing compressionWriter", err.Error())
				}
			}(compressionWriter)
			defer func(compressionWriter *gzip.Writer) {
				err := compressionWriter.Close()
				if err != nil {
					logger.ZapSugarLogger.Errorln("error when closing compressionWriter", err.Error())
				}
			}(compressionWriter)

			response.Header().Set("Content-Encoding", "gzip")

			next.ServeHTTP(
				// here we override simple writer with compression writer
				gzipWriter{ResponseWriter: response, Writer: compressionWriter},
				request,
			)
			return
		}
		logger.ZapSugarLogger.Debugln("gzip not accepted for this request")
		next.ServeHTTP(response, request)
	})
}

func getCompressedWriter(response http.ResponseWriter) *gzip.Writer {
	compressionWriter, err := gzip.NewWriterLevel(response, gzip.BestSpeed)
	if err != nil {
		_, err := io.WriteString(response, err.Error())
		if err != nil {
			logger.ZapSugarLogger.Errorln("error when creation gzip writer ", err.Error())
		}
		return nil
	}

	return compressionWriter
}

func decompressRequest(request *http.Request) *http.Request {
	if request == nil {
		return request
	}

	var bodyBuf bytes.Buffer
	_, err := request.Body.Read(bodyBuf.Bytes())
	if err != nil && err.Error() == "EOF" {
		return request
	}
	if err != nil {
		logger.ZapSugarLogger.Errorln("error when reading body ", err.Error())
		return request
	}
	logger.ZapSugarLogger.Debugln("decompressing body ", bodyBuf.String())

	compressionReader, err := gzip.NewReader(request.Body)
	if err != nil {
		logger.ZapSugarLogger.Errorln("error when creating gzip reader ", err.Error())

		return request
	}
	request.Body = compressionReader

	return request
}

func isGzipAvailableForThisRequest(request *http.Request) (isOk bool) {
	if request == nil {
		return false
	}

	return isCorrectAcceptEncoding(request)
}

func isRequestCompressed(request *http.Request) (isOk bool) {
	if request == nil {
		return false
	}

	return isContentEncodingGzip(request)
}

func isContentEncodingGzip(request *http.Request) bool {
	contentEncoding := request.Header.Values("Content-Encoding")

	for i := 0; i < len(contentEncoding); i++ {
		if strings.Contains(contentEncoding[i], "gzip") {
			return true
		}
	}
	return false
}

func isCorrectAcceptEncoding(request *http.Request) bool {
	acceptEncodings := request.Header.Values("Accept-Encoding")
	for i := 0; i < len(acceptEncodings); i++ {
		if strings.Contains(acceptEncodings[i], "gzip") {
			return true
		}
	}
	return false
}
