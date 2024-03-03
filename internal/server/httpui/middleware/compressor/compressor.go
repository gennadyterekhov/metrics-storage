package compressor

import (
	"compress/gzip"
	"github.com/gennadyterekhov/metrics-storage/internal/constants"
	"github.com/gennadyterekhov/metrics-storage/internal/logger"
	"io"
	"net/http"
	"strings"
)

type decompressedBody struct {
	io.ReadCloser
	Reader io.Reader
}

type gzipWriter struct {
	http.ResponseWriter
	// gzip.NewWriterLevel will be here
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func (dcb decompressedBody) Read(b []byte) (int, error) {
	return dcb.Reader.Read(b)
}

func GzipCompressor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request != nil && IsGzipAvailableForThisRequest(request) {
			compressionWriter, err := gzip.NewWriterLevel(response, gzip.BestSpeed)
			if err != nil {
				_, err := io.WriteString(response, err.Error())
				if err != nil {
					logger.ZapSugarLogger.Warnln("error when creation gzip writer ", err.Error())
				}
				return
			}
			defer func(compressionWriter *gzip.Writer) {
				err := compressionWriter.Flush()
				if err != nil {
					logger.ZapSugarLogger.Warnln("error when flushing compressionWriter", err.Error())
				}
			}(compressionWriter)

			defer func(compressionWriter *gzip.Writer) {
				err := compressionWriter.Close()
				if err != nil {
					logger.ZapSugarLogger.Warnln("error when closing compressionWriter", err.Error())
				}
			}(compressionWriter)

			response.Header().Set("Content-Encoding", "gzip")

			//decoded by default
			//decompressedBodyReader, err := getDecompressedBodyReader(request)
			//if err != nil {
			//	logger.ZapSugarLogger.Debugln("error when getting decompressed body reader", err.Error())
			//	_, err := io.WriteString(response, err.Error())
			//	if err != nil {
			//		logger.ZapSugarLogger.Debugln("error when writing error to http response", err.Error())
			//	}
			//	return
			//}
			//
			//request.Body = decompressedBody{
			//	Reader: decompressedBodyReader,
			//	//BodyBytes: bodyBytes,
			//}
			next.ServeHTTP(
				// here we override simple writer with compression writer
				gzipWriter{ResponseWriter: response, Writer: compressionWriter},
				request,
			)
			response.Header().Set("Content-Encoding", "gzip")
			return
		} else {
			logger.ZapSugarLogger.Debugln("gzip not accepted for this request")
		}

		next.ServeHTTP(response, request)
	})
}

func IsGzipAvailableForThisRequest(request *http.Request) (isOk bool) {
	correctAcceptContentType := false
	correctAcceptEncoding := false

	acceptContentTypes := request.Header.Values("Accept")
	acceptEncodings := request.Header.Values("Accept-Encoding")
	contentEncodings := request.Header.Values("Content-Encoding")

	for i := 0; i < len(contentEncodings); i += 1 {
		if strings.Contains(contentEncodings[i], "gzip") {
			// if body is already encoded no further checks are necessary
			//correctContentEncoding = true
			return true
		}
	}
	for i := 0; i < len(acceptContentTypes); i += 1 {
		if strings.Contains(acceptContentTypes[i], constants.TextHTML) ||
			strings.Contains(acceptContentTypes[i], "html/text") ||
			strings.Contains(acceptContentTypes[i], constants.ApplicationJSON) {
			correctAcceptContentType = true
			break
		}
	}
	for i := 0; i < len(acceptEncodings); i += 1 {
		if strings.Contains(acceptEncodings[i], "gzip") {
			correctAcceptEncoding = true
			break
		}
	}

	return correctAcceptContentType && correctAcceptEncoding
}

func getDecompressedBodyReader(r *http.Request) (gz *gzip.Reader, err error) {
	// создаём *gzip.Reader, который будет читать тело запроса
	// и распаковывать его
	gz, err = gzip.NewReader(r.Body)
	if err != nil {
		logger.ZapSugarLogger.Warnln("error when opening gzip body reader", err.Error())

		return nil, err
	}

	return gz, nil
}
