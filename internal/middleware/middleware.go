package middleware

import (
	"github.com/gennadyterekhov/metrics-storage/internal/exceptions"
	"github.com/gennadyterekhov/metrics-storage/internal/services"
	"github.com/gennadyterekhov/metrics-storage/internal/validators"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Middleware func(http.Handler) http.Handler

func Conveyor(h http.Handler, middlewares ...Middleware) http.Handler {
	for _, middlewareCallback := range middlewares {
		h = middlewareCallback(h)
	}
	return h
}

func MethodPost(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			http.Error(res, exceptions.UpdateMetricsMethodNotAllowed, http.StatusMethodNotAllowed)
			return
		}
		next.ServeHTTP(res, req)
	})
}

func MethodGet(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			http.Error(res, exceptions.GetOneMetricsMethodNotAllowed, http.StatusMethodNotAllowed)
			return
		}
		next.ServeHTTP(res, req)
	})
}

func URLParametersToGetMetricsArePresent(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		metricType, name, err := validators.GetDataToGet(
			chi.URLParam(req, "metricType"),
			chi.URLParam(req, "metricName"),
		)
		if err != nil && err.Error() == exceptions.InvalidMetricTypeChoice {
			http.Error(res, err.Error(), http.StatusNotFound)
			return
		}
		if err != nil && err.Error() == exceptions.InvalidMetricType {
			http.Error(res, err.Error(), http.StatusNotImplemented)
			return
		}
		if err != nil && err.Error() == exceptions.EmptyMetricName {
			http.Error(res, err.Error(), http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		_, err = services.GetMetricAsString(metricType, name)

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

		next.ServeHTTP(res, req)

	})
}

func URLParametersToSetMetricsArePresent(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		_, _, _, _, err := validators.GetDataToSave(
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

		next.ServeHTTP(res, req)
	})
}
