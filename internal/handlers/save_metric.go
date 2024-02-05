package handlers

import (
	"github.com/gennadyterekhov/metrics-storage/internal/services"
	"github.com/gennadyterekhov/metrics-storage/internal/validators"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func SaveMetric(res http.ResponseWriter, req *http.Request) {

	metricType, name, counterValue, gaugeValue, err := validators.GetDataToSave(
		chi.URLParam(req, "metricType"),
		chi.URLParam(req, "metricName"),
		chi.URLParam(req, "metricValue"),
	)

	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	services.SaveMetricToMemory(metricType, name, counterValue, gaugeValue)
}
