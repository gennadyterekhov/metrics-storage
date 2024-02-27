package handlers

import (
	"encoding/json"
	"github.com/gennadyterekhov/metrics-storage/internal/constants"
	"github.com/gennadyterekhov/metrics-storage/internal/constants/exceptions"
	"github.com/gennadyterekhov/metrics-storage/internal/domain/models"
	"github.com/gennadyterekhov/metrics-storage/internal/server/app/services/get_metric_service"
	"github.com/gennadyterekhov/metrics-storage/internal/server/httpui/validators"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
)

func GetMetric(res http.ResponseWriter, req *http.Request) {
	if req.Header.Get(constants.HeaderContentType) == constants.ApplicationJSON {
		decoder := json.NewDecoder(req.Body)
		metrics := models.Metrics{}
		if err := decoder.Decode(&metrics); err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		metricType := metrics.MType
		name := metrics.ID

		metric, err := getmetricservice.GetMetricsAsStruct(metricType, name)
		//if metricType == types.Counter {
		//	metrics.Delta += metric.Delta
		//
		//} else {
		//	metrics.Value = metric
		//}
		if err != nil {
			http.Error(res, err.Error(), http.StatusNotFound)
			return
		}
		res.Header().Set(constants.HeaderContentType, constants.ApplicationJSON)

		encoder := json.NewEncoder(res)
		if err := encoder.Encode(metric); err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		return
	}

	metricType, name, err := validators.GetDataToGet(
		chi.URLParam(req, "metricType"),
		chi.URLParam(req, "metricName"),
	)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	metric, err := getmetricservice.GetMetricAsString(metricType, name)

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
