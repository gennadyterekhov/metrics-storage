package router

import (
	"github.com/gennadyterekhov/metrics-storage/internal/common"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			resp, get := common.SendRequest(t, ts, http.MethodPost, tt.url)
			assert.Equal(t, tt.want.code, resp.StatusCode)
			assert.Equal(t, tt.want.response, get)
		})
	}
}
