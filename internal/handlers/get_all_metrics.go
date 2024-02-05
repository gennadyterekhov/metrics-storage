package handlers

import (
	"github.com/gennadyterekhov/metrics-storage/internal/services"
	"io"
	"net/http"
)

func GetAllMetrics(res http.ResponseWriter, req *http.Request) {

	htmlPage := services.GetMetricsListAsHTML()

	_, err := io.WriteString(res, htmlPage)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}
