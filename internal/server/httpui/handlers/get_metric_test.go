package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gennadyterekhov/metrics-storage/internal/constants"
	"github.com/gennadyterekhov/metrics-storage/internal/constants/types"
	"github.com/gennadyterekhov/metrics-storage/internal/container"
	"github.com/gennadyterekhov/metrics-storage/internal/domain/models"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

type args struct {
	typ  string
	name string
}

func TestGetMetricJSON(t *testing.T) {
	container.MetricsRepository.Clear()
	container.MetricsRepository.AddCounter("cnt", 1)

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
			request := httptest.NewRequest(http.MethodPost, "/value", body)
			request.Header.Set(constants.HeaderContentType, constants.ApplicationJSON)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("metricType", tt.args.typ)
			rctx.URLParams.Add("metricName", tt.args.name)
			request.Header.Set(constants.HeaderContentType, constants.ApplicationJSON)
			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))
			w := httptest.NewRecorder()
			GetMetricJSONHandler()(w, request)

			res := w.Result()
			assert.Equal(t, tt.want.code, res.StatusCode)

			defer res.Body.Close()
			if res.StatusCode == http.StatusOK {

				metricFromResponse, _ := io.ReadAll(res.Body)
				receivedData := models.Metrics{}
				_ = json.Unmarshal(metricFromResponse, &receivedData)
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
	container.MetricsRepository.Clear()
	container.MetricsRepository.AddCounter("cnt", 1)

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
			request := httptest.NewRequest(http.MethodGet, url, nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("metricType", tt.args.typ)
			rctx.URLParams.Add("metricName", tt.args.name)
			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))
			w := httptest.NewRecorder()
			GetMetricHandler()(w, request)

			res := w.Result()
			metricFromResponse, _ := io.ReadAll(res.Body)
			metricFromResponseAsInt, _ := strconv.ParseInt(string(metricFromResponse), 10, 64)
			defer res.Body.Close()
			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, tt.want.metricValue, metricFromResponseAsInt)

		})
	}
}
