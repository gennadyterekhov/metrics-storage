package logger

import (
	"github.com/gennadyterekhov/metrics-storage/internal/constants"
	"github.com/gennadyterekhov/metrics-storage/internal/logger"
	"net/http"
	"time"
)

type (
	LogContext struct {
		uri       string
		method    string
		startTime time.Time
		time      time.Duration
		status    int
		size      int
	}

	LoggingResponseWriter struct {
		http.ResponseWriter
		LogContext *LogContext
	}
)

func (lrw *LoggingResponseWriter) Write(b []byte) (int, error) {
	size, err := lrw.ResponseWriter.Write(b)
	lrw.LogContext.size += size
	return size, err
}

func (lrw *LoggingResponseWriter) WriteHeader(statusCode int) {
	logger.ZapSugarLogger.Debugln("writing header from log middleware")
	lrw.ResponseWriter.WriteHeader(statusCode)
	lrw.LogContext.status = statusCode
}

func (lrw *LoggingResponseWriter) log() {
	logger.ZapSugarLogger.Infoln(
		"uri", lrw.LogContext.uri,
		"method", lrw.LogContext.method,
		"duration", lrw.LogContext.time,
		"status", lrw.LogContext.status,
		"size", lrw.LogContext.size,
	)
}
func (lrw *LoggingResponseWriter) updateContext(req *http.Request) {
	if req != nil {
		lrw.LogContext.uri = req.RequestURI
		lrw.LogContext.method = req.Method
		lrw.LogContext.time = time.Since(lrw.LogContext.startTime)
	} else {
		logger.ZapSugarLogger.Debugln("could not update log ctx with actual request info")
	}
}

func RequestAndResponseLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req == nil {
			logger.ZapSugarLogger.Debugln("logger middleware: request is nil")
			next.ServeHTTP(res, req)
			return
		}
		reqBody := make([]byte, 0)
		_, _ = req.Body.Read(reqBody)
		// don't close because it will be used in compressor middleware
		//defer req.Body.Close()
		logger.ZapSugarLogger.Debugln(
			"got request",
			req.Method,
			req.RequestURI,
			req.Header.Get(constants.HeaderContentType),
			string(reqBody),
		)

		customWriter := initializeCustomWriter(res, req)
		if customWriter == nil {
			logger.ZapSugarLogger.Debugln("could not set custom logger writer")
		}
		next.ServeHTTP(customWriter, req)

		if customWriter != nil {
			customWriter.updateContext(req)
			customWriter.log()
		}
	})
}

func initializeCustomWriter(res http.ResponseWriter, req *http.Request) *LoggingResponseWriter {
	responseData := &LogContext{
		uri:       req.RequestURI,
		method:    req.Method,
		time:      0,
		startTime: time.Now(),
	}
	customWriter := LoggingResponseWriter{
		ResponseWriter: res,
		LogContext:     responseData,
	}
	return &customWriter
}
