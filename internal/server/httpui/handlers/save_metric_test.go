package handlers

import (
	"bytes"
	"context"
	"github.com/gennadyterekhov/metrics-storage/internal/constants"
	"github.com/gennadyterekhov/metrics-storage/internal/constants/types"
	"github.com/gennadyterekhov/metrics-storage/internal/container"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSaveMetricHttpMethodJSON(t *testing.T) {
	type want struct {
		code        int
		response    string
		typ         string
		metricName  string
		metricValue int64
	}
	tests := []struct {
		name   string
		method string
		want   want
	}{
		{
			name:   "ok",
			method: http.MethodPost,
			want:   want{code: http.StatusOK},
		},
		{
			name:   "-",
			method: http.MethodGet,
			want:   want{code: http.StatusMethodNotAllowed},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rawJSON := `{"id":"cnt", "type":"counter", "delta":1}`
			//request := httptest.NewRequest(tt.method, "/update/counter/cnt/1", bytes.NewBuffer([]byte(rawJSON)))
			request := httptest.NewRequest(tt.method, "/update", bytes.NewBuffer([]byte(rawJSON)))

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("metricType", types.Counter)
			rctx.URLParams.Add("metricName", "cnt")
			rctx.URLParams.Add("metricValue", "1")
			request.Header.Set(constants.HeaderContentType, constants.ApplicationJSON)
			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))
			w := httptest.NewRecorder()
			SaveMetricHandler()(w, request)

			res := w.Result()
			defer res.Body.Close()
			assert.Equal(t, tt.want.code, res.StatusCode)
		})
	}
}

func TestSaveMetricJSON(t *testing.T) {
	type want struct {
		code        int
		response    string
		typ         string
		metricName  string
		metricValue int64
	}
	tests := []struct {
		name    string
		url     string
		rawJSON string
		want    want
	}{
		{
			name:    "Counter",
			url:     "/update/counter/cnt/1",
			rawJSON: `{"id":"cnt", "type":"counter", "delta":1}`,
			want:    want{code: http.StatusOK, response: "", typ: types.Counter, metricName: "cnt", metricValue: 1},
		},
		{
			name:    "Gauge",
			url:     "/update/gauge/gaugeName/1",
			rawJSON: `{"id":"gaugeName", "type":"gauge", "value":1}`,
			want:    want{code: http.StatusOK, response: "", typ: types.Gauge, metricName: "gaugeName", metricValue: 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			container.MetricsRepository.Clear()
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("metricType", tt.want.typ)
			rctx.URLParams.Add("metricName", tt.want.metricName)
			rctx.URLParams.Add("metricValue", "1")
			request := httptest.NewRequest(http.MethodPost, tt.url, bytes.NewBuffer([]byte(tt.rawJSON)))
			request.Header.Set(constants.HeaderContentType, constants.ApplicationJSON)
			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

			w := httptest.NewRecorder()
			SaveMetricHandler()(w, request)

			res := w.Result()
			defer res.Body.Close()
			assert.Equal(t, tt.want.code, res.StatusCode)

			if tt.want.typ == types.Counter {
				assert.Equal(t, tt.want.metricValue, container.MetricsRepository.GetCounterOrZero(tt.want.metricName))
			}
			if tt.want.typ == types.Gauge {
				assert.Equal(t, tt.want.metricValue, int64(container.MetricsRepository.GetGaugeOrZero(tt.want.metricName)))
			}
		})
	}

	// check counter is added to itself
	container.MetricsRepository.AddCounter("cnt", 1)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("metricType", types.Counter)
	rctx.URLParams.Add("metricName", "cnt")
	rctx.URLParams.Add("metricValue", "10")
	request := httptest.NewRequest(http.MethodPost, "/update/counter/cnt/10", bytes.NewBuffer([]byte(`{"id":"cnt", "type":"counter", "delta":10}`)))
	request.Header.Set(constants.HeaderContentType, constants.ApplicationJSON)
	request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()
	SaveMetricHandler()(w, request)
	assert.Equal(t, int64(10+1), container.MetricsRepository.GetCounterOrZero("cnt"))

	// check gauge is substituted
	container.MetricsRepository.SetGauge("gaugeName", 1)

	rctx = chi.NewRouteContext()
	rctx.URLParams.Add("metricType", types.Gauge)
	rctx.URLParams.Add("metricName", "gaugeName")
	rctx.URLParams.Add("metricValue", "3")
	request = httptest.NewRequest(http.MethodPost, "/update/gauge/gaugeName/3", bytes.NewBuffer([]byte(`{"id":"gaugeName", "type":"gauge", "value":3}`)))
	request.Header.Set(constants.HeaderContentType, constants.ApplicationJSON)
	request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

	w = httptest.NewRecorder()
	SaveMetricHandler()(w, request)
	assert.Equal(t, float64(3), container.MetricsRepository.GetGaugeOrZero("gaugeName"))
}

func TestSaveMetricHttpMethod(t *testing.T) {
	type want struct {
		code        int
		response    string
		typ         string
		metricName  string
		metricValue int64
	}
	tests := []struct {
		name   string
		method string
		want   want
	}{
		{
			name:   "ok",
			method: http.MethodPost,
			want:   want{code: http.StatusOK},
		},
		{
			name:   "-",
			method: http.MethodGet,
			want:   want{code: http.StatusMethodNotAllowed},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, "/update/counter/cnt/1", nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("metricType", types.Counter)
			rctx.URLParams.Add("metricName", "cnt")
			rctx.URLParams.Add("metricValue", "1")

			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))
			w := httptest.NewRecorder()
			SaveMetricHandler()(w, request)

			res := w.Result()
			defer res.Body.Close()
			assert.Equal(t, tt.want.code, res.StatusCode)
		})
	}
}

func TestSaveMetric(t *testing.T) {
	type want struct {
		code        int
		response    string
		typ         string
		metricName  string
		metricValue int64
	}
	tests := []struct {
		name string
		url  string
		want want
	}{
		{
			name: "Counter",
			url:  "/update/counter/cnt/1",
			want: want{code: http.StatusOK, response: "", typ: types.Counter, metricName: "cnt", metricValue: 1},
		},
		{
			name: "Gauge",
			url:  "/update/gauge/gaugeName/1",
			want: want{code: http.StatusOK, response: "", typ: types.Gauge, metricName: "gaugeName", metricValue: 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			container.MetricsRepository.Clear()
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("metricType", tt.want.typ)
			rctx.URLParams.Add("metricName", tt.want.metricName)
			rctx.URLParams.Add("metricValue", "1")
			request := httptest.NewRequest(http.MethodPost, tt.url, nil)
			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

			w := httptest.NewRecorder()
			SaveMetricHandler()(w, request)

			res := w.Result()
			defer res.Body.Close()
			assert.Equal(t, tt.want.code, res.StatusCode)

			if tt.want.typ == types.Counter {
				assert.Equal(t, tt.want.metricValue, container.MetricsRepository.GetCounterOrZero(tt.want.metricName))
			}
			if tt.want.typ == types.Gauge {
				assert.Equal(t, tt.want.metricValue, int64(container.MetricsRepository.GetGaugeOrZero(tt.want.metricName)))
			}
		})
	}

	// check counter is added to itself
	container.MetricsRepository.AddCounter("cnt", 1)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("metricType", types.Counter)
	rctx.URLParams.Add("metricName", "cnt")
	rctx.URLParams.Add("metricValue", "10")
	request := httptest.NewRequest(http.MethodPost, "/update/counter/cnt/10", nil)
	request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()
	SaveMetricHandler()(w, request)
	assert.Equal(t, int64(10+1), container.MetricsRepository.GetCounterOrZero("cnt"))

	// check gauge is substituted
	container.MetricsRepository.SetGauge("gaugeName", 1)

	rctx = chi.NewRouteContext()
	rctx.URLParams.Add("metricType", types.Gauge)
	rctx.URLParams.Add("metricName", "gaugeName")
	rctx.URLParams.Add("metricValue", "3")
	request = httptest.NewRequest(http.MethodPost, "/update/gauge/gaugeName/3", nil)
	request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

	w = httptest.NewRecorder()
	SaveMetricHandler()(w, request)
	assert.Equal(t, float64(3), container.MetricsRepository.GetGaugeOrZero("gaugeName"))
}