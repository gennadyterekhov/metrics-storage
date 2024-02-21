package logger

import (
	"go.uber.org/zap"
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
	lrw.LogContext.status = statusCode
}

func (lrw *LoggingResponseWriter) log() {
	ZapSugarLogger.Infoln(
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
		ZapSugarLogger.Errorln("could not update log context with actual request info")
	}
}

var ZapSugarLogger zap.SugaredLogger

func RequestAndResponseLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		customWriter := initializeCustomWriter(res, req)
		if customWriter == nil {
			ZapSugarLogger.Errorln("could not set custom logger writer")
		}
		next.ServeHTTP(res, req)

		if customWriter != nil {
			customWriter.updateContext(req)
			customWriter.log()
		}
	})
}

func initializeCustomWriter(res http.ResponseWriter, req *http.Request) *LoggingResponseWriter {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	// делаем регистратор SugaredLogger
	ZapSugarLogger = *logger.Sugar()

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
