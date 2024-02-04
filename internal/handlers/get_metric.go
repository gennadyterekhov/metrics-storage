package handlers

import (
	"github.com/gennadyterekhov/metrics-storage/internal/exceptions"
	"github.com/gennadyterekhov/metrics-storage/internal/services"
	"github.com/gennadyterekhov/metrics-storage/internal/validators"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
)

func GetMetric(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, exceptions.GetOneMetricsMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	metricType, name, err := validators.GetDataToGet(
		chi.URLParam(req, "metricType"),
		chi.URLParam(req, "metricName"),
	)
	if err != nil && err.Error() == exceptions.InvalidMetricTypeChoice {
		http.Error(res, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil && err.Error() == exceptions.EmptyMetricName {
		http.Error(res, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil && err.Error() == exceptions.InvalidMetricType {
		http.Error(res, err.Error(), http.StatusNotImplemented)
		return
	}
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	metric := services.GetMetricAsString(metricType, name)

	_, err = io.WriteString(res, metric)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}
