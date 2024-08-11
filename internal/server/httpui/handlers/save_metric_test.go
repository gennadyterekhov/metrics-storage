package handlers

import (
	"bytes"
	"context"
	"net/http"
	"testing"

	"github.com/gennadyterekhov/metrics-storage/internal/common/tests"
	"github.com/stretchr/testify/suite"

	"github.com/gennadyterekhov/metrics-storage/internal/common/constants/types"
	"github.com/gennadyterekhov/metrics-storage/internal/common/logger"
	"github.com/gennadyterekhov/metrics-storage/internal/common/testhelper"
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

func (suite *saveMetricTestSuite) TestSaveMetricHttpMethodJSON() {
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
		suite.T().Run(tt.name, func(t *testing.T) {
			rawJSON := `{"id":"cnt", "type":"counter", "delta":1}`
			response, _ := testhelper.SendAlreadyJSONedBody(
				t,
				suite.TestHTTPServer.Server,
				tt.method,
				"/update/counter/cnt/1",
				bytes.NewBuffer([]byte(rawJSON)),
			)
			err := response.Body.Close()
			if err != nil {
				return
			}

			assert.Equal(suite.T(), tt.want.code, response.StatusCode)
		})
	}
}

func (suite *saveMetricTestSuite) TestSaveMetricJSON() {
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
		suite.T().Run(tt.name, func(t *testing.T) {
			suite.Repository.Clear()

			response, _ := testhelper.SendAlreadyJSONedBody(
				t,
				suite.TestHTTPServer.Server,
				http.MethodPost,
				tt.url,
				bytes.NewBuffer([]byte(tt.rawJSON)),
			)
			err := response.Body.Close()
			if err != nil {
				return
			}

			assert.Equal(suite.T(), tt.want.code, response.StatusCode)

			if tt.want.typ == types.Counter {
				assert.Equal(suite.T(), tt.want.metricValue, suite.Repository.GetCounterOrZero(context.Background(), tt.want.metricName))
			}
			if tt.want.typ == types.Gauge {
				assert.Equal(suite.T(), tt.want.metricValue, int64(suite.Repository.GetGaugeOrZero(context.Background(), tt.want.metricName)))
			}
		})
	}

	// check counter is added to itself
	suite.Repository.AddCounter(context.Background(), "cnt", 1)
	response, _ := testhelper.SendAlreadyJSONedBody(
		suite.T(),
		suite.TestHTTPServer.Server,
		http.MethodPost,
		"/update/counter/cnt/10",
		bytes.NewBuffer([]byte(`{"id":"cnt", "type":"counter", "delta":10}`)),
	)
	err := response.Body.Close()
	if err != nil {
		return
	}

	assert.Equal(suite.T(), int64(10+1), suite.Repository.GetCounterOrZero(context.Background(), "cnt"))

	// check gauge is substituted
	suite.Repository.SetGauge(context.Background(), "gaugeName", 1)
	response, _ = testhelper.SendAlreadyJSONedBody(
		suite.T(),
		suite.TestHTTPServer.Server,
		http.MethodPost,
		"/update/gauge/gaugeName/3",
		bytes.NewBuffer([]byte(`{"id":"gaugeName", "type":"gauge", "value":3}`)),
	)
	err = response.Body.Close()
	if err != nil {
		return
	}

	assert.Equal(suite.T(), float64(3), suite.Repository.GetGaugeOrZero(context.Background(), "gaugeName"))
}

func (suite *saveMetricTestSuite) TestSaveMetricJSONReturnsUpdatedValuesInBody() {
	rawJSON := `{"id":"cnt", "type":"counter", "delta":10}`
	suite.Repository.Clear()
	suite.Repository.AddCounter(context.Background(), "cnt", 1)

	response, responseBody := testhelper.SendAlreadyJSONedBody(
		suite.T(),
		suite.TestHTTPServer.Server,
		http.MethodPost,
		"/update/counter/cnt/10",
		bytes.NewBuffer([]byte(rawJSON)),
	)
	err := response.Body.Close()
	if err != nil {
		return
	}

	assert.Equal(suite.T(), http.StatusOK, response.StatusCode)

	logger.Custom.Debugln("responseBody", responseBody)

	assert.Equal(suite.T(), `{"type":"counter","id":"cnt","delta":11}`, string(responseBody))
}

func (suite *saveMetricTestSuite) TestSaveMetricHttpMethod() {
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
		suite.T().Run(tt.name, func(t *testing.T) {
			response, _ := testhelper.SendRequest(
				t,
				suite.TestHTTPServer.Server,
				tt.method,
				"/update/counter/cnt/1",
			)
			err := response.Body.Close()
			if err != nil {
				return
			}

			assert.Equal(suite.T(), tt.want.code, response.StatusCode)
		})
	}
}

func (suite *saveMetricTestSuite) TestSaveMetric() {
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
		suite.T().Run(tt.name, func(t *testing.T) {
			suite.Repository.Clear()
			response, _ := testhelper.SendRequest(
				t,
				suite.TestHTTPServer.Server,
				http.MethodPost,
				tt.url,
			)
			err := response.Body.Close()
			if err != nil {
				return
			}

			assert.Equal(suite.T(), tt.want.code, response.StatusCode)

			if tt.want.typ == types.Counter {
				assert.Equal(suite.T(), tt.want.metricValue, suite.Repository.GetCounterOrZero(context.Background(), tt.want.metricName))
			}
			if tt.want.typ == types.Gauge {
				assert.Equal(suite.T(), tt.want.metricValue, int64(suite.Repository.GetGaugeOrZero(context.Background(), tt.want.metricName)))
			}
		})
	}

	// check counter is added to itself
	suite.Repository.AddCounter(context.Background(), "cnt", 1)
	response, _ := testhelper.SendRequest(
		suite.T(),
		suite.TestHTTPServer.Server,
		http.MethodPost,
		"/update/counter/cnt/10",
	)
	err := response.Body.Close()
	if err != nil {
		return
	}

	assert.Equal(suite.T(), int64(10+1), suite.Repository.GetCounterOrZero(context.Background(), "cnt"))

	// check gauge is substituted
	suite.Repository.SetGauge(context.Background(), "gaugeName", 1)
	response, _ = testhelper.SendRequest(
		suite.T(),
		suite.TestHTTPServer.Server,
		http.MethodPost,
		"/update/gauge/gaugeName/3",
	)
	err = response.Body.Close()
	if err != nil {
		return
	}

	assert.Equal(suite.T(), float64(3), suite.Repository.GetGaugeOrZero(context.Background(), "gaugeName"))
}

func (suite *saveMetricTestSuite) TestGzipCompression() {
	requestBody := `{"id":"cnt", "type":"counter", "delta":1}`
	successBody := `{"id":"cnt", "type":"counter", "delta":1}`

	suite.T().Run("client can send gzipped request", func(t *testing.T) {
		suite.Repository.Clear()

		response, _ := testhelper.SendGzipRequest(
			t,
			suite.TestHTTPServer.Server,
			http.MethodPost,
			"/update/",
			requestBody,
		)
		err := response.Body.Close()
		if err != nil {
			return
		}

		require.Equal(t, http.StatusOK, response.StatusCode)
	})

	suite.T().Run("client can send gzipped request and server can respond with gzipped body", func(t *testing.T) {
		suite.Repository.Clear()

		suite.Repository.AddCounter(context.Background(), "cnt", 1)

		response, responseBody := testhelper.SendGzipRequest(
			t,
			suite.TestHTTPServer.Server,
			http.MethodPost,
			"/value",
			requestBody,
		)
		err := response.Body.Close()
		if err != nil {
			return
		}
		require.Equal(t, http.StatusOK, response.StatusCode)
		require.JSONEq(t, successBody, string(responseBody))
	})
}

func (suite *saveMetricTestSuite) TestCanSaveMetricToDB() {
	suite.T().Skip("only manual use because depends on host")

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
		suite.T().Run(tt.name, func(t *testing.T) {
			suite.Repository.Clear()
			response, _ := testhelper.SendRequest(
				t,
				suite.TestHTTPServer.Server,
				http.MethodPost,
				tt.url,
			)
			err := response.Body.Close()
			if err != nil {
				return
			}

			assert.Equal(suite.T(), tt.want.code, response.StatusCode)

			if tt.want.typ == types.Counter {
				assert.Equal(suite.T(), tt.want.metricValue, suite.Repository.GetCounterOrZero(context.Background(), tt.want.metricName))
			}
			if tt.want.typ == types.Gauge {
				assert.Equal(suite.T(), tt.want.metricValue, int64(suite.Repository.GetGaugeOrZero(context.Background(), tt.want.metricName)))
			}
		})
	}

	// check counter is added to itself
	suite.Repository.AddCounter(context.Background(), "cnt", 1)
	response, _ := testhelper.SendRequest(
		suite.T(),
		suite.TestHTTPServer.Server,
		http.MethodPost,
		"/update/counter/cnt/10",
	)
	err := response.Body.Close()
	if err != nil {
		return
	}

	assert.Equal(suite.T(), int64(10+1), suite.Repository.GetCounterOrZero(context.Background(), "cnt"))

	// check gauge is substituted
	suite.Repository.SetGauge(context.Background(), "gaugeName", 1)
	response, _ = testhelper.SendRequest(
		suite.T(),
		suite.TestHTTPServer.Server,
		http.MethodPost,
		"/update/gauge/gaugeName/3",
	)
	err = response.Body.Close()
	if err != nil {
		return
	}

	assert.Equal(suite.T(), float64(3), suite.Repository.GetGaugeOrZero(context.Background(), "gaugeName"))
}

func (suite *saveMetricTestSuite) TestSaveMetricList() {
	suite.Repository.Clear()
	rawJSON := `[
					{"id":"Alloc", "type":"gauge", "value":1.1},
					{"id":"BuckHashSys", "type":"gauge", "value":2.2},
					{"id":"PollCount", "type":"counter", "delta":3}
	]`
	suite.T().Run("save list", func(t *testing.T) {
		response, _ := testhelper.SendAlreadyJSONedBody(
			t,
			suite.TestHTTPServer.Server,
			http.MethodPost,
			"/updates/",
			bytes.NewBuffer([]byte(rawJSON)),
		)
		err := response.Body.Close()
		if err != nil {
			return
		}

		assert.Equal(suite.T(), http.StatusOK, response.StatusCode)

		assert.Equal(suite.T(), 2, len(suite.Repository.GetAllGauges(context.Background())))
		assert.Equal(suite.T(), 1, len(suite.Repository.GetAllCounters(context.Background())))
	})
}
