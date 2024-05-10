package handlers

import (
	"io"
	"net/http"

	"github.com/gennadyterekhov/metrics-storage/internal/constants"
	"github.com/gennadyterekhov/metrics-storage/internal/server/app"
	"github.com/gennadyterekhov/metrics-storage/internal/server/httpui/middleware"
)

func GetAllMetrics(res http.ResponseWriter, req *http.Request) {
	htmlPage := app.GetMetricsListAsHTML(req.Context())

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
