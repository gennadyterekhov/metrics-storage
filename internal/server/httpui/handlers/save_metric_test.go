package handlers

import (
	"bytes"
	"github.com/gennadyterekhov/metrics-storage/internal/constants/types"
	"github.com/gennadyterekhov/metrics-storage/internal/logger"
	"github.com/gennadyterekhov/metrics-storage/internal/server/storage"
	"github.com/gennadyterekhov/metrics-storage/internal/testhelper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
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
			response, _ := testhelper.SendAlreadyJSONedBody(
				t,
				testhelper.TestServer,
				tt.method,
				"/update/counter/cnt/1",
				bytes.NewBuffer([]byte(rawJSON)),
			)
			response.Body.Close()

			assert.Equal(t, tt.want.code, response.StatusCode)
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
			storage.MetricsRepository.Clear()

			response, _ := testhelper.SendAlreadyJSONedBody(
				t,
				testhelper.TestServer,
				http.MethodPost,
				tt.url,
				bytes.NewBuffer([]byte(tt.rawJSON)),
			)
			response.Body.Close()

			assert.Equal(t, tt.want.code, response.StatusCode)

			if tt.want.typ == types.Counter {
				assert.Equal(t, tt.want.metricValue, storage.MetricsRepository.GetCounterOrZero(tt.want.metricName))
			}
			if tt.want.typ == types.Gauge {
				assert.Equal(t, tt.want.metricValue, int64(storage.MetricsRepository.GetGaugeOrZero(tt.want.metricName)))
			}
		})
	}

	// check counter is added to itself
	storage.MetricsRepository.AddCounter("cnt", 1)
	response, _ := testhelper.SendAlreadyJSONedBody(
		t,
		testhelper.TestServer,
		http.MethodPost,
		"/update/counter/cnt/10",
		bytes.NewBuffer([]byte(`{"id":"cnt", "type":"counter", "delta":10}`)),
	)
	response.Body.Close()

	assert.Equal(t, int64(10+1), storage.MetricsRepository.GetCounterOrZero("cnt"))

	// check gauge is substituted
	storage.MetricsRepository.SetGauge("gaugeName", 1)
	response, _ = testhelper.SendAlreadyJSONedBody(
		t,
		testhelper.TestServer,
		http.MethodPost,
		"/update/gauge/gaugeName/3",
		bytes.NewBuffer([]byte(`{"id":"gaugeName", "type":"gauge", "value":3}`)),
	)
	response.Body.Close()

	assert.Equal(t, float64(3), storage.MetricsRepository.GetGaugeOrZero("gaugeName"))
}

func TestSaveMetricJSONReturnsUpdatedValuesInBody(t *testing.T) {
	rawJSON := `{"id":"cnt", "type":"counter", "delta":10}`
	storage.MetricsRepository.Clear()
	storage.MetricsRepository.AddCounter("cnt", 1)

	response, responseBody := testhelper.SendAlreadyJSONedBody(
		t,
		testhelper.TestServer,
		http.MethodPost,
		"/update/counter/cnt/10",
		bytes.NewBuffer([]byte(rawJSON)),
	)
	response.Body.Close()

	assert.Equal(t, http.StatusOK, response.StatusCode)

	logger.ZapSugarLogger.Debugln("responseBody", responseBody)

	assert.Equal(t, `{"type":"counter","id":"cnt","delta":11}`, string(responseBody))
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
			response, _ := testhelper.SendRequest(
				t,
				testhelper.TestServer,
				tt.method,
				"/update/counter/cnt/1",
			)
			response.Body.Close()

			assert.Equal(t, tt.want.code, response.StatusCode)
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
		{
			name: "invalid_value",
			url:  "/update/counter/testCounter/none",
			want: want{code: http.StatusBadRequest, response: "", typ: types.Counter, metricName: "testCounter", metricValue: 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage.MetricsRepository.Clear()
			response, _ := testhelper.SendRequest(
				t,
				testhelper.TestServer,
				http.MethodPost,
				tt.url,
			)
			response.Body.Close()

			assert.Equal(t, tt.want.code, response.StatusCode)

			if tt.want.typ == types.Counter {
				assert.Equal(t, tt.want.metricValue, storage.MetricsRepository.GetCounterOrZero(tt.want.metricName))
			}
			if tt.want.typ == types.Gauge {
				assert.Equal(t, tt.want.metricValue, int64(storage.MetricsRepository.GetGaugeOrZero(tt.want.metricName)))
			}
		})
	}

	// check counter is added to itself
	storage.MetricsRepository.AddCounter("cnt", 1)
	response, _ := testhelper.SendRequest(
		t,
		testhelper.TestServer,
		http.MethodPost,
		"/update/counter/cnt/10",
	)
	response.Body.Close()

	assert.Equal(t, int64(10+1), storage.MetricsRepository.GetCounterOrZero("cnt"))

	// check gauge is substituted
	storage.MetricsRepository.SetGauge("gaugeName", 1)
	response, _ = testhelper.SendRequest(
		t,
		testhelper.TestServer,
		http.MethodPost,
		"/update/gauge/gaugeName/3",
	)
	response.Body.Close()

	assert.Equal(t, float64(3), storage.MetricsRepository.GetGaugeOrZero("gaugeName"))
}

func TestGzipCompression(t *testing.T) {
	requestBody := `{"id":"cnt", "type":"counter", "delta":1}`
	successBody := `{"id":"cnt", "type":"counter", "delta":1}`

	t.Run("client can send gzipped request", func(t *testing.T) {
		storage.MetricsRepository.Clear()

		response, _ := testhelper.SendGzipRequest(
			t,
			testhelper.TestServer,
			http.MethodPost,
			"/update/",
			requestBody,
		)
		response.Body.Close()

		require.Equal(t, http.StatusOK, response.StatusCode)
	})

	t.Run("client can send gzipped request and server can respond with gzipped body", func(t *testing.T) {
		storage.MetricsRepository.Clear()

		storage.MetricsRepository.AddCounter("cnt", 1)

		response, responseBody := testhelper.SendGzipRequest(
			t,
			testhelper.TestServer,
			http.MethodPost,
			"/value",
			requestBody,
		)
		response.Body.Close()
		require.Equal(t, http.StatusOK, response.StatusCode)
		require.JSONEq(t, successBody, string(responseBody))
	})
}
