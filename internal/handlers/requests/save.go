package requests

import (
	"fmt"
	"github.com/gennadyterekhov/metrics-storage/internal/exceptions"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type RawSaveMetricRequest struct {
	MetricType  string
	MetricName  string
	MetricValue string
}

type ValidatedSaveCounterMetricRequest ValidatedSaveMetricRequest[int64]
type ValidatedSaveGaugeMetricRequest ValidatedSaveMetricRequest[float64]

type ValidatedSaveMetricRequest[T int64 | float64] struct {
	MetricType  string
	MetricName  string
	MetricValue T
}

func GetRawRequest(res http.ResponseWriter, req *http.Request) (*RawSaveMetricRequest, error) {
	metricType := chi.URLParam(req, "metricType")
	metricName := chi.URLParam(req, "metricName")
	metricValue := chi.URLParam(req, "metricValue")
	if metricType == "" {
		http.Error(res, exceptions.EmptyMetricType, http.StatusBadRequest)
		return nil, fmt.Errorf(exceptions.EmptyMetricType)
	}
	if metricName == "" {
		http.Error(res, exceptions.EmptyMetricName, http.StatusBadRequest)
		return nil, fmt.Errorf(exceptions.EmptyMetricName)
	}
	if metricValue == "" {
		http.Error(res, exceptions.EmptyMetricValue, http.StatusBadRequest)
		return nil, fmt.Errorf(exceptions.EmptyMetricValue)
	}

	return &RawSaveMetricRequest{
		MetricType:  metricType,
		MetricName:  metricName,
		MetricValue: metricValue,
	}, nil
}
