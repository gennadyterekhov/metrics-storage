package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/gennadyterekhov/metrics-storage/internal/common/constants"
	"github.com/gennadyterekhov/metrics-storage/internal/common/constants/types"
	"github.com/gennadyterekhov/metrics-storage/internal/common/domain/models"
	"github.com/gennadyterekhov/metrics-storage/internal/common/testhelper"
	"github.com/gennadyterekhov/metrics-storage/internal/server/config"
	"github.com/gennadyterekhov/metrics-storage/internal/server/storage"
	"github.com/stretchr/testify/assert"
)

type args struct {
	typ  string
	name string
}

func TestGetMetricJSON(t *testing.T) {
	storage.MetricsRepository.AddCounter(context.Background(), "cnt", 1)

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
		t.Run(tt.name, func(t *testing.T) {
			body := getBodyFromArgs(tt.args)
			response, responseBody := testhelper.SendAlreadyJSONedBody(
				t,
				testhelper.TestServer,
				http.MethodPost,
				"/value",
				body,
			)
			response.Body.Close()
			assert.Equal(t, tt.want.code, response.StatusCode)

			if response.StatusCode == http.StatusOK {
				receivedData := models.Metrics{}
				err := json.Unmarshal(responseBody, &receivedData)
				assert.NoError(t, err)
				assert.Equal(t, tt.want.metricValue, *receivedData.Delta)
			}
		})
	}
}

func getBodyFromArgs(arguments args) *bytes.Buffer {
	rawJSON := fmt.Sprintf(`{"id":"%s", "type":"%s"}`, arguments.name, arguments.typ)

	return bytes.NewBuffer([]byte(rawJSON))
}

func TestGetMetric(t *testing.T) {
	storage.MetricsRepository.Clear()
	storage.MetricsRepository.AddCounter(context.Background(), "cnt", 1)

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
		t.Run(tt.name, func(t *testing.T) {
			url := "/value/" + tt.args.typ + "/" + tt.args.name

			response, responseBody := testhelper.SendRequest(
				t,
				testhelper.TestServer,
				http.MethodGet,
				url,
			)
			response.Body.Close()

			metricFromResponseAsInt, _ := strconv.ParseInt(string(responseBody), 10, 64)
			assert.Equal(t, tt.want.code, response.StatusCode)
			assert.Equal(t, tt.want.metricValue, metricFromResponseAsInt)
		})
	}
}

func TestCanGetMetricFromDB(t *testing.T) {
	t.Skip("only manual use because depends on host")

	config.Conf.DBDsn = constants.TestDBDsn
	storage.MetricsRepository = storage.CreateDBStorage()

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
		t.Run(tt.name, func(t *testing.T) {
			storage.MetricsRepository.Clear()
			if tt.want.code == http.StatusOK {
				storage.MetricsRepository.AddCounter(context.Background(), "cnt", 1)
			}

			url := "/value/" + tt.args.typ + "/" + tt.args.name

			response, responseBody := testhelper.SendRequest(
				t,
				testhelper.TestServer,
				http.MethodGet,
				url,
			)
			response.Body.Close()

			metricFromResponseAsInt, _ := strconv.ParseInt(string(responseBody), 10, 64)
			assert.Equal(t, tt.want.code, response.StatusCode)
			assert.Equal(t, tt.want.metricValue, metricFromResponseAsInt)
		})
	}
	config.Conf.DBDsn = ""
	storage.MetricsRepository.CloseDB()
	storage.MetricsRepository = storage.CreateRAMStorage()
}
