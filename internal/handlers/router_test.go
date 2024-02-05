package handlers

import (
	"github.com/gennadyterekhov/metrics-storage/internal/common"
	"github.com/gennadyterekhov/metrics-storage/internal/exceptions"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCanGetUrlParameters(t *testing.T) {
	ts := httptest.NewServer(GetRouter())
	defer ts.Close()

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
			want: want{code: http.StatusOK, response: "", typ: "counter", metricName: "cnt", metricValue: 1},
		},
		{
			name: "Gauge",
			url:  "/update/gauge/gaugeName/1",
			want: want{code: http.StatusOK, response: "", typ: "gauge", metricName: "gaugeName", metricValue: 1},
		},
		{
			name: "wrong type status code",
			url:  "/update/unknown/testCounter/100",
			want: want{code: http.StatusBadRequest, response: exceptions.InvalidMetricTypeChoice, typ: "gauge", metricName: "gaugeName", metricValue: 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			resp, get := common.SendRequest(t, ts, http.MethodPost, tt.url)
			defer resp.Body.Close()
			assert.Equal(t, tt.want.code, resp.StatusCode)
			assert.Equal(t, tt.want.response, strings.Trim(get, "\n"))
		})
	}
}
