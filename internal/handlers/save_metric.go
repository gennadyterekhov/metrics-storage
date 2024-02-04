package handlers

import (
	"github.com/gennadyterekhov/metrics-storage/internal/exceptions"
	"github.com/gennadyterekhov/metrics-storage/internal/services"
	"github.com/gennadyterekhov/metrics-storage/internal/validators"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func SaveMetric(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, exceptions.UpdateMetricsMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	metricType, name, counterValue, gaugeValue, err := validators.GetDataToSave(
		chi.URLParam(req, "metricType"),
		chi.URLParam(req, "metricName"),
		chi.URLParam(req, "metricValue"),
	)
	if err != nil && err.Error() == exceptions.InvalidMetricTypeChoice {
		http.Error(res, err.Error(), http.StatusBadRequest)
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

	services.SaveMetricToMemory(metricType, name, counterValue, gaugeValue)
}
