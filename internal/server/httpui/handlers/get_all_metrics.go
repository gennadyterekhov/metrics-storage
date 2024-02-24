package handlers

import (
	"github.com/gennadyterekhov/metrics-storage/internal/server/app/services/get_metric_service"
	"io"
	"net/http"
)

func GetAllMetrics(res http.ResponseWriter, req *http.Request) {

	htmlPage := getmetricservice.GetMetricsListAsHTML()

	_, err := io.WriteString(res, htmlPage)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}
