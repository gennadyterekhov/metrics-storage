package handlers

import (
	"bytes"
	"context"
	"net/http"
	"testing"

	"github.com/gennadyterekhov/metrics-storage/internal/common/tests"
	"github.com/stretchr/testify/suite"

	"github.com/gennadyterekhov/metrics-storage/internal/common/constants"
	"github.com/gennadyterekhov/metrics-storage/internal/common/constants/types"
	"github.com/gennadyterekhov/metrics-storage/internal/common/logger"
	"github.com/gennadyterekhov/metrics-storage/internal/common/testhelper"
	"github.com/gennadyterekhov/metrics-storage/internal/server/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type saveMetricTestSuite struct {
	tests.BaseSuiteWithServer
}

func (suite *saveMetricTestSuite) SetupSuite() {
	tests.InitBaseSuiteWithServer(suite)
}

func TestSaveMetricHandler(t *testing.T) {
	suite.Run(t, new(saveMetricTestSuite))
}

func (st *saveMetricTestSuite) TestSaveMetricHttpMethodJSON() {
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
		st.T().Run(tt.name, func(t *testing.T) {
			rawJSON := `{"id":"cnt", "type":"counter", "delta":1}`
			response, _ := testhelper.SendAlreadyJSONedBody(
				t,
				st.TestHTTPServer.Server,
				tt.method,
				"/update/counter/cnt/1",
				bytes.NewBuffer([]byte(rawJSON)),
			)
			response.Body.Close()

			assert.Equal(st.T(), tt.want.code, response.StatusCode)
		})
	}
}

func (st *saveMetricTestSuite) TestSaveMetricJSON() {
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
		st.T().Run(tt.name, func(t *testing.T) {
			st.Repository.Clear()

			response, _ := testhelper.SendAlreadyJSONedBody(
				t,
				st.TestHTTPServer.Server,
				http.MethodPost,
				tt.url,
				bytes.NewBuffer([]byte(tt.rawJSON)),
			)
			response.Body.Close()

			assert.Equal(st.T(), tt.want.code, response.StatusCode)

			if tt.want.typ == types.Counter {
				assert.Equal(st.T(), tt.want.metricValue, st.Repository.GetCounterOrZero(context.Background(), tt.want.metricName))
			}
			if tt.want.typ == types.Gauge {
				assert.Equal(st.T(), tt.want.metricValue, int64(st.Repository.GetGaugeOrZero(context.Background(), tt.want.metricName)))
			}
		})
	}

	// check counter is added to itself
	st.Repository.AddCounter(context.Background(), "cnt", 1)
	response, _ := testhelper.SendAlreadyJSONedBody(
		st.T(),
		st.TestHTTPServer.Server,
		http.MethodPost,
		"/update/counter/cnt/10",
		bytes.NewBuffer([]byte(`{"id":"cnt", "type":"counter", "delta":10}`)),
	)
	response.Body.Close()

	assert.Equal(st.T(), int64(10+1), st.Repository.GetCounterOrZero(context.Background(), "cnt"))

	// check gauge is substituted
	st.Repository.SetGauge(context.Background(), "gaugeName", 1)
	response, _ = testhelper.SendAlreadyJSONedBody(
		st.T(),
		st.TestHTTPServer.Server,
		http.MethodPost,
		"/update/gauge/gaugeName/3",
		bytes.NewBuffer([]byte(`{"id":"gaugeName", "type":"gauge", "value":3}`)),
	)
	response.Body.Close()

	assert.Equal(st.T(), float64(3), st.Repository.GetGaugeOrZero(context.Background(), "gaugeName"))
}

func (st *saveMetricTestSuite) TestSaveMetricJSONReturnsUpdatedValuesInBody() {
	rawJSON := `{"id":"cnt", "type":"counter", "delta":10}`
	st.Repository.Clear()
	st.Repository.AddCounter(context.Background(), "cnt", 1)

	response, responseBody := testhelper.SendAlreadyJSONedBody(
		st.T(),
		st.TestHTTPServer.Server,
		http.MethodPost,
		"/update/counter/cnt/10",
		bytes.NewBuffer([]byte(rawJSON)),
	)
	response.Body.Close()

	assert.Equal(st.T(), http.StatusOK, response.StatusCode)

	logger.ZapSugarLogger.Debugln("responseBody", responseBody)

	assert.Equal(st.T(), `{"type":"counter","id":"cnt","delta":11}`, string(responseBody))
}

func (st *saveMetricTestSuite) TestSaveMetricHttpMethod() {
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
		st.T().Run(tt.name, func(t *testing.T) {
			response, _ := testhelper.SendRequest(
				t,
				st.TestHTTPServer.Server,
				tt.method,
				"/update/counter/cnt/1",
			)
			response.Body.Close()

			assert.Equal(st.T(), tt.want.code, response.StatusCode)
		})
	}
}

func (st *saveMetricTestSuite) TestSaveMetric() {
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
		st.T().Run(tt.name, func(t *testing.T) {
			st.Repository.Clear()
			response, _ := testhelper.SendRequest(
				t,
				st.TestHTTPServer.Server,
				http.MethodPost,
				tt.url,
			)
			response.Body.Close()

			assert.Equal(st.T(), tt.want.code, response.StatusCode)

			if tt.want.typ == types.Counter {
				assert.Equal(st.T(), tt.want.metricValue, st.Repository.GetCounterOrZero(context.Background(), tt.want.metricName))
			}
			if tt.want.typ == types.Gauge {
				assert.Equal(st.T(), tt.want.metricValue, int64(st.Repository.GetGaugeOrZero(context.Background(), tt.want.metricName)))
			}
		})
	}

	// check counter is added to itself
	st.Repository.AddCounter(context.Background(), "cnt", 1)
	response, _ := testhelper.SendRequest(
		st.T(),
		st.TestHTTPServer.Server,
		http.MethodPost,
		"/update/counter/cnt/10",
	)
	response.Body.Close()

	assert.Equal(st.T(), int64(10+1), st.Repository.GetCounterOrZero(context.Background(), "cnt"))

	// check gauge is substituted
	st.Repository.SetGauge(context.Background(), "gaugeName", 1)
	response, _ = testhelper.SendRequest(
		st.T(),
		st.TestHTTPServer.Server,
		http.MethodPost,
		"/update/gauge/gaugeName/3",
	)
	response.Body.Close()

	assert.Equal(st.T(), float64(3), st.Repository.GetGaugeOrZero(context.Background(), "gaugeName"))
}

func (st *saveMetricTestSuite) TestGzipCompression() {
	requestBody := `{"id":"cnt", "type":"counter", "delta":1}`
	successBody := `{"id":"cnt", "type":"counter", "delta":1}`

	st.T().Run("client can send gzipped request", func(t *testing.T) {
		st.Repository.Clear()

		response, _ := testhelper.SendGzipRequest(
			t,
			st.TestHTTPServer.Server,
			http.MethodPost,
			"/update/",
			requestBody,
		)
		response.Body.Close()

		require.Equal(t, http.StatusOK, response.StatusCode)
	})

	st.T().Run("client can send gzipped request and server can respond with gzipped body", func(t *testing.T) {
		st.Repository.Clear()

		st.Repository.AddCounter(context.Background(), "cnt", 1)

		response, responseBody := testhelper.SendGzipRequest(
			t,
			st.TestHTTPServer.Server,
			http.MethodPost,
			"/value",
			requestBody,
		)
		response.Body.Close()
		require.Equal(t, http.StatusOK, response.StatusCode)
		require.JSONEq(t, successBody, string(responseBody))
	})
}

func (st *saveMetricTestSuite) TestCanSaveMetricToDB() {
	st.T().Skip("only manual use because depends on host")

	config.Conf.DBDsn = constants.TestDBDsn
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
		st.T().Run(tt.name, func(t *testing.T) {
			st.Repository.Clear()
			response, _ := testhelper.SendRequest(
				t,
				st.TestHTTPServer.Server,
				http.MethodPost,
				tt.url,
			)
			response.Body.Close()

			assert.Equal(st.T(), tt.want.code, response.StatusCode)

			if tt.want.typ == types.Counter {
				assert.Equal(st.T(), tt.want.metricValue, st.Repository.GetCounterOrZero(context.Background(), tt.want.metricName))
			}
			if tt.want.typ == types.Gauge {
				assert.Equal(st.T(), tt.want.metricValue, int64(st.Repository.GetGaugeOrZero(context.Background(), tt.want.metricName)))
			}
		})
	}

	// check counter is added to itself
	st.Repository.AddCounter(context.Background(), "cnt", 1)
	response, _ := testhelper.SendRequest(
		st.T(),
		st.TestHTTPServer.Server,
		http.MethodPost,
		"/update/counter/cnt/10",
	)
	response.Body.Close()

	assert.Equal(st.T(), int64(10+1), st.Repository.GetCounterOrZero(context.Background(), "cnt"))

	// check gauge is substituted
	st.Repository.SetGauge(context.Background(), "gaugeName", 1)
	response, _ = testhelper.SendRequest(
		st.T(),
		st.TestHTTPServer.Server,
		http.MethodPost,
		"/update/gauge/gaugeName/3",
	)
	response.Body.Close()

	assert.Equal(st.T(), float64(3), st.Repository.GetGaugeOrZero(context.Background(), "gaugeName"))
	config.Conf.DBDsn = ""
}

func (st *saveMetricTestSuite) TestSaveMetricList() {
	st.Repository.Clear()
	rawJSON := `[
					{"id":"Alloc", "type":"gauge", "value":1.1},
					{"id":"BuckHashSys", "type":"gauge", "value":2.2},
					{"id":"PollCount", "type":"counter", "delta":3}
	]`
	st.T().Run("save list", func(t *testing.T) {
		response, _ := testhelper.SendAlreadyJSONedBody(
			t,
			st.TestHTTPServer.Server,
			http.MethodPost,
			"/updates/",
			bytes.NewBuffer([]byte(rawJSON)),
		)
		response.Body.Close()

		assert.Equal(st.T(), http.StatusOK, response.StatusCode)

		assert.Equal(st.T(), 2, len(st.Repository.GetAllGauges(context.Background())))
		assert.Equal(st.T(), 1, len(st.Repository.GetAllCounters(context.Background())))
	})
}
