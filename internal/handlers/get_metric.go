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
	metricType, name, err := validators.GetDataToGet(
		chi.URLParam(req, "metricType"),
		chi.URLParam(req, "metricName"),
	)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	metric, err := services.GetMetricAsString(metricType, name)

	if err != nil && err.Error() == exceptions.UnknownMetricName {
		http.Error(res, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil && err.Error() == exceptions.InvalidMetricTypeChoice {
		http.Error(res, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = io.WriteString(res, metric)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}
