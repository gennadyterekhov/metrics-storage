package logger

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/gennadyterekhov/metrics-storage/internal/common/constants"
	"github.com/gennadyterekhov/metrics-storage/internal/common/logger"
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
	logger.Custom.Debugln("writing header from log middleware, status:", statusCode)
	lrw.ResponseWriter.WriteHeader(statusCode)
	lrw.LogContext.status = statusCode
}

func (lrw *LoggingResponseWriter) log() {
	logger.Custom.Infoln(
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
		logger.Custom.Debugln("could not update log ctx with actual request info")
	}
}

func RequestAndResponseLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req == nil {
			logger.Custom.Debugln("logger middleware: request is nil")
			next.ServeHTTP(res, req)
			return
		}
		var reqBody []byte
		reqBody, err := io.ReadAll(req.Body)
		if err != nil {
			logger.Custom.Errorln("could not read body", err.Error())
		}
		req.Body = io.NopCloser(bytes.NewBuffer(reqBody))

		logger.Custom.Debugln(
			"got request",
			req.Method,
			req.RequestURI,
			req.Header.Get(constants.HeaderContentType),
			string(reqBody),
		)

		customWriter := initializeCustomWriter(res, req)
		if customWriter == nil {
			logger.Custom.Debugln("could not set custom logger writer")
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
