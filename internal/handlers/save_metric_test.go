package handlers

import (
	"errors"
	"github.com/gennadyterekhov/metrics-storage/internal/container"
	"github.com/gennadyterekhov/metrics-storage/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestSaveMetricHttpMethod(t *testing.T) {
	t.Skipf("cannot get url parameters from chi when running method without router")

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
			w := httptest.NewRecorder()
			SaveMetric(w, request)

			res := w.Result()
			defer res.Body.Close()
			assert.Equal(t, tt.want.code, res.StatusCode)
		})
	}
}

func TestSaveMetric(t *testing.T) {
	t.Skipf("cannot get url parameters from chi when running method without router")
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
			request := httptest.NewRequest(http.MethodPost, tt.url, nil)
			w := httptest.NewRecorder()
			SaveMetric(w, request)

			res := w.Result()
			defer res.Body.Close()
			assert.Equal(t, tt.want.code, res.StatusCode)

			if tt.want.typ == types.Counter {
				assert.Equal(t, tt.want.metricValue, container.Instance.MetricsRepository.GetCounterOrZero(tt.want.metricName))
			}
			if tt.want.typ == types.Gauge {
				assert.Equal(t, tt.want.metricValue, int64(container.Instance.MetricsRepository.GetGaugeOrZero(tt.want.metricName)))
			}
		})
	}

	// check counter is added to itself
	request := httptest.NewRequest(http.MethodPost, "/update/counter/cnt/10", nil)
	w := httptest.NewRecorder()
	SaveMetric(w, request)
	assert.Equal(t, int64(10+1), container.Instance.MetricsRepository.GetCounterOrZero("cnt"))

	// check gauge is substituted
	request = httptest.NewRequest(http.MethodPost, "/update/gauge/gaugeName/3", nil)
	w = httptest.NewRecorder()
	SaveMetric(w, request)
	assert.Equal(t, float64(3), container.Instance.MetricsRepository.GetGaugeOrZero("gaugeName"))
}

func Test_parseURL(t *testing.T) {
	type args struct {
		url string
	}
	possibleError := errors.New("expected exactly 3 parameters")
	tests := []struct {
		name           string
		args           args
		wantMetricType string
		wantName       string
		wantValue      string
		wantErr        bool
	}{
		{
			name:           "ok",
			args:           args{url: "/update/counter/cnt/1"},
			wantMetricType: "counter",
			wantName:       "cnt",
			wantValue:      "1",
			wantErr:        false,
		},
		{
			name:           "too short 1",
			args:           args{url: "/update/"},
			wantMetricType: "counter",
			wantName:       "cnt",
			wantValue:      "1",
			wantErr:        true,
		},
		{
			name:           "too short 2",
			args:           args{url: "/update/counter"},
			wantMetricType: "counter",
			wantName:       "cnt",
			wantValue:      "1",
			wantErr:        true,
		},
		{
			name:           "too short 3",
			args:           args{url: "/update/counter/cnt"},
			wantMetricType: "counter",
			wantName:       "cnt",
			wantValue:      "1",
			wantErr:        true,
		},
		{
			name:           "too long",
			args:           args{url: "/update/counter/cnt/1/hello"},
			wantMetricType: "counter",
			wantName:       "cnt",
			wantValue:      "1",
			wantErr:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			urlPath := &url.URL{Path: tt.args.url}

			gotMetricType, gotName, gotValue, err := parseURL(urlPath)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			if (err != nil) != tt.wantErr {
				if err != nil {
					require.NoError(t, err)
				}

				t.Errorf("parseURL() error = %v, wantErr %v", err, possibleError)
				return
			}
			if gotMetricType != tt.wantMetricType {
				t.Errorf("parseURL() gotMetricType = %v, want %v", gotMetricType, tt.wantMetricType)
			}
			if gotName != tt.wantName {
				t.Errorf("parseURL() gotName = %v, want %v", gotName, tt.wantName)
			}
			if gotValue != tt.wantValue {
				t.Errorf("parseURL() gotValue = %v, want %v", gotValue, tt.wantValue)
			}
		})
	}
}
