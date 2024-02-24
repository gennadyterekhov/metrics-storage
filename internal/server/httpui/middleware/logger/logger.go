package logger

import (
	"fmt"
	"github.com/gennadyterekhov/metrics-storage/internal/constants"
	"github.com/gennadyterekhov/metrics-storage/internal/logger"
	"net/http"
	"time"
)

type (
	LogContext struct {
		status    int
		size      int
		uri       string
		method    string
		startTime time.Time
		time      time.Duration
	}

	LoggingResponseWriter struct /* implements http.ResponseWriter*/ {
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
	lrw.ResponseWriter.WriteHeader(statusCode)
	fmt.Println("raw print statusCode", statusCode)
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
		logger.ZapSugarLogger.Errorln("could not update log context with actual request info")
	}
}

func RequestAndResponseLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		reqBody := make([]byte, 0)
		_, _ = req.Body.Read(reqBody)
		defer req.Body.Close()
		logger.ZapSugarLogger.Debugln(
			"got request",
			req.Method,
			req.RequestURI,
			req.Header.Get(constants.HeaderContentType),
			reqBody,
		)

		customWriter := initializeCustomWriter(res, req)
		if customWriter == nil {
			logger.ZapSugarLogger.Errorln("could not set custom logger writer")
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
