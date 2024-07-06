package handlers

import (
	"io"
	"net/http"

	"github.com/gennadyterekhov/metrics-storage/internal/server/services"

	"github.com/gennadyterekhov/metrics-storage/internal/common/constants"
	"github.com/gennadyterekhov/metrics-storage/internal/server/httpui/middleware"
)

func GetAllMetrics(res http.ResponseWriter, req *http.Request) {
	htmlPage := services.GetMetricsListAsHTML(req.Context())

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
