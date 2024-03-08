package handlers

import (
	"github.com/gennadyterekhov/metrics-storage/internal/constants"
	"github.com/gennadyterekhov/metrics-storage/internal/server/app"
	"github.com/gennadyterekhov/metrics-storage/internal/server/httpui/middleware"
	"io"
	"net/http"
)

func GetAllMetrics(res http.ResponseWriter, req *http.Request) {

	htmlPage := app.GetMetricsListAsHTML()

	res.Header().Set(constants.HeaderContentType, constants.TextHTML)
	_, err := io.WriteString(res, htmlPage)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}
func GetAllMetricsHandler() http.Handler {
	return middleware.CommonConveyor(
		http.HandlerFunc(GetAllMetrics),
	)
}
func GetAllMetricsHandlerFunc() func(http.ResponseWriter, *http.Request) {
	return GetAllMetricsHandler().ServeHTTP
}
