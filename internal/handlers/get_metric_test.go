package handlers

import (
	"github.com/gennadyterekhov/metrics-storage/internal/container"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestGetMetric(t *testing.T) {
	container.Instance.MetricsRepository.AddCounter("cnt", 1)
	type want struct {
		code        int
		metricValue int64
	}
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "",
			args: args{url: "/value/counter/cnt"},
			want: want{code: http.StatusOK, metricValue: 1},
		},
		{
			name: "",
			args: args{url: "/value/counter/unknown"},
			want: want{code: http.StatusNotFound, metricValue: 0},
		},
		{
			name: "",
			args: args{url: "/value/unknown/cnt"},
			want: want{code: http.StatusNotFound, metricValue: 0},
		},
		{
			name: "",
			args: args{url: "/value/counter/"},
			want: want{code: http.StatusNotFound, metricValue: 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			request := httptest.NewRequest(http.MethodGet, tt.args.url, nil)
			w := httptest.NewRecorder()
			GetMetric(w, request)

			res := w.Result()
			metricFromResponse, _ := io.ReadAll(res.Body)
			metricFromResponseAsInt, _ := strconv.ParseInt(string(metricFromResponse), 10, 64)
			defer res.Body.Close()
			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, tt.want.metricValue, metricFromResponseAsInt)

		})
	}
}
