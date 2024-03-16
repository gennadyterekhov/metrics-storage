package compressor

import (
	"compress/gzip"
	"github.com/gennadyterekhov/metrics-storage/internal/constants"
	"github.com/gennadyterekhov/metrics-storage/internal/logger"
	"io"
	"net/http"
	"strings"
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

		if IsGzipAvailableForThisRequest(request) {
			compressionWriter := getCompressedWriter(response)
			if compressionWriter == nil {
				return
			}
			defer compressionWriter.Flush()
			defer compressionWriter.Close()

			response.Header().Set("Content-Encoding", "gzip")
			response.Header().Set(constants.HeaderContentType, constants.ApplicationJSON)

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
			logger.ZapSugarLogger.Warnln("error when creation gzip writer ", err.Error())
		}
		return nil
	}

	return compressionWriter
}

func IsGzipAvailableForThisRequest(request *http.Request) (isOk bool) {
	if request == nil {
		return false
	}

	return isCorrectAcceptEncoding(request) && isCorrectContentType(request)

}

func isCorrectContentType(request *http.Request) bool {
	correctContentType := false
	contentType := request.Header.Get("Content-Type")
	correctContentType = contentType == constants.ApplicationJSON ||
		contentType == constants.TextHTML ||
		contentType == "html/text"
	return correctContentType
}

func isCorrectAcceptEncoding(request *http.Request) bool {
	correctAcceptEncoding := false
	acceptEncodings := request.Header.Values("Accept-Encoding")
	for i := 0; i < len(acceptEncodings); i += 1 {
		if strings.Contains(acceptEncodings[i], "gzip") {
			correctAcceptEncoding = true
			break
		}
	}
	return correctAcceptEncoding
}
