package handlers

import (
	"encoding/json"
	"github.com/gennadyterekhov/metrics-storage/internal/constants"
	"github.com/gennadyterekhov/metrics-storage/internal/constants/types"
	"github.com/gennadyterekhov/metrics-storage/internal/domain/dto"
	"github.com/gennadyterekhov/metrics-storage/internal/domain/models"
	"github.com/gennadyterekhov/metrics-storage/internal/logger"
	"github.com/gennadyterekhov/metrics-storage/internal/server/app/services/save_metric_service"
	"github.com/gennadyterekhov/metrics-storage/internal/server/httpui/validators"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func SaveMetric(res http.ResponseWriter, req *http.Request) {
	if req.Header.Get(constants.HeaderContentType) == constants.ApplicationJSON {

		decoder := json.NewDecoder(req.Body)
		metric := models.Metrics{}

		if err := decoder.Decode(&metric); err != nil {
			logger.ZapSugarLogger.Debugln("could not decode json body", err.Error())

			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		filledDto := dto.MetricToSaveDto{}
		if metric.MType == types.Counter {
			filledDto = dto.MetricToSaveDto{
				Type:         metric.MType,
				Name:         metric.ID,
				CounterValue: *metric.Delta,
			}
		} else {

			filledDto = dto.MetricToSaveDto{
				Type:       metric.MType,
				Name:       metric.ID,
				GaugeValue: *metric.Value,
			}
		}

		res.Header().Set(constants.HeaderContentType, constants.ApplicationJSON)
		savemetricservice.SaveMetricToMemory(&filledDto)

		encoder := json.NewEncoder(res)
		if err := encoder.Encode(metric); err != nil {
			logger.ZapSugarLogger.Debugln("could not encode json body", err.Error())
			return
		}
		return
	}
	filledDto, err := validators.GetDataToSave(
		chi.URLParam(req, "metricType"),
		chi.URLParam(req, "metricName"),
		chi.URLParam(req, "metricValue"),
	)

	if err != nil {
		logger.ZapSugarLogger.Debugln("validation error", err.Error())
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	savemetricservice.SaveMetricToMemory(filledDto)
}
