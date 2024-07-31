package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/gennadyterekhov/metrics-storage/internal/server/httpui/responses"

	"github.com/gennadyterekhov/metrics-storage/internal/common/tests"
	"github.com/stretchr/testify/suite"

	"github.com/gennadyterekhov/metrics-storage/internal/common/constants/types"
	"github.com/gennadyterekhov/metrics-storage/internal/common/testhelper"
	"github.com/stretchr/testify/assert"
)

type args struct {
	typ  string
	name string
}

type getMetricTestSuite struct {
	tests.BaseSuiteWithServer
}

func (st *getMetricTestSuite) SetupSuite() {
	tests.InitBaseSuiteWithServer(st)
}

func TestGetMetricHandler(t *testing.T) {
	suite.Run(t, new(getMetricTestSuite))
}

func (st *getMetricTestSuite) TestGetMetricJSON() {
	st.Repository.AddCounter(context.Background(), "cnt", 1)

	type want struct {
		code        int
		metricValue int64
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "ok",
			args: args{typ: types.Counter, name: "cnt"},
			want: want{code: http.StatusOK, metricValue: 1},
		},
		{
			name: "name unknown",
			args: args{typ: types.Counter, name: "unknown"},
			want: want{code: http.StatusNotFound, metricValue: 0},
		},
		{
			name: "type unknown",
			args: args{typ: "unknown", name: "cnt"},
			want: want{code: http.StatusNotFound, metricValue: 0},
		},
		{
			name: "name empty",
			args: args{typ: types.Counter, name: ""},
			want: want{code: http.StatusNotFound, metricValue: 0},
		},
	}
	for _, tt := range tests {
		st.T().Run(tt.name, func(t *testing.T) {
			body := getBodyFromArgs(tt.args)
			response, responseBody := testhelper.SendAlreadyJSONedBody(
				t,
				st.TestHTTPServer.Server,
				http.MethodPost,
				"/value",
				body,
			)
			err := response.Body.Close()
			if err != nil {
				return
			}
			assert.Equal(t, tt.want.code, response.StatusCode)

			if response.StatusCode == http.StatusOK {
				receivedData := responses.GetMetricResponse{}
				err := json.Unmarshal(responseBody, &receivedData)
				assert.NoError(t, err)
				assert.Equal(t, tt.want.metricValue, *receivedData.CounterValue)
			}
		})
	}
}

func getBodyFromArgs(arguments args) *bytes.Buffer {
	rawJSON := fmt.Sprintf(`{"id":"%s", "type":"%s"}`, arguments.name, arguments.typ)

	return bytes.NewBuffer([]byte(rawJSON))
}

func (st *getMetricTestSuite) TestGetMetric() {
	st.Repository.AddCounter(context.Background(), "cnt", 1)

	type want struct {
		code        int
		metricValue int64
	}
	type args struct {
		typ  string
		name string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "ok",
			args: args{typ: types.Counter, name: "cnt"},
			want: want{code: http.StatusOK, metricValue: 1},
		},
		{
			name: "name unknown",
			args: args{typ: types.Counter, name: "unknown"},
			want: want{code: http.StatusNotFound, metricValue: 0},
		},
		{
			name: "type unknown",
			args: args{typ: "unknown", name: "cnt"},
			want: want{code: http.StatusNotFound, metricValue: 0},
		},
		{
			name: "name empty",
			args: args{typ: types.Counter, name: ""},
			want: want{code: http.StatusNotFound, metricValue: 0},
		},
	}
	for _, tt := range tests {
		st.T().Run(tt.name, func(t *testing.T) {
			url := "/value/" + tt.args.typ + "/" + tt.args.name

			response, responseBody := testhelper.SendRequest(
				t,
				st.TestHTTPServer.Server,
				http.MethodGet,
				url,
			)
			err := response.Body.Close()
			assert.NoError(t, err)

			metricFromResponseAsInt, err := strconv.ParseInt(string(responseBody), 10, 64) //
			if tt.name == "ok" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.want.code, response.StatusCode)
			assert.Equal(t, tt.want.metricValue, metricFromResponseAsInt)
		})
	}
}

func (st *getMetricTestSuite) TestCanGetMetricFromDB() {
	st.T().Skip("only manual use because depends on host")
	type want struct {
		code        int
		metricValue int64
	}
	type args struct {
		typ  string
		name string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "ok",
			args: args{typ: types.Counter, name: "cnt"},
			want: want{code: http.StatusOK, metricValue: 1},
		},
		{
			name: "name unknown",
			args: args{typ: types.Counter, name: "unknown"},
			want: want{code: http.StatusNotFound, metricValue: 0},
		},
		{
			name: "type unknown",
			args: args{typ: "unknown", name: "cnt"},
			want: want{code: http.StatusNotFound, metricValue: 0},
		},
		{
			name: "name empty",
			args: args{typ: types.Counter, name: ""},
			want: want{code: http.StatusNotFound, metricValue: 0},
		},
	}
	for _, tt := range tests {
		st.T().Run(tt.name, func(t *testing.T) {
			if tt.want.code == http.StatusOK {
				st.Repository.AddCounter(context.Background(), "cnt", 1)
			}

			url := "/value/" + tt.args.typ + "/" + tt.args.name

			response, responseBody := testhelper.SendRequest(
				t,
				st.TestHTTPServer.Server,
				http.MethodGet,
				url,
			)
			err := response.Body.Close()
			assert.NoError(t, err)

			metricFromResponseAsInt, err := strconv.ParseInt(string(responseBody), 10, 64)
			if tt.name == "ok" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.want.code, response.StatusCode)
			assert.Equal(t, tt.want.metricValue, metricFromResponseAsInt)
		})
	}
}
