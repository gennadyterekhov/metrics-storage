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
				assert.Equal(t, tt.want.metricValue, container.Instance.MemStorage.Counters[tt.want.metricName])
			}
			if tt.want.typ == types.Gauge {
				assert.Equal(t, tt.want.metricValue, int64(container.Instance.MemStorage.Gauges[tt.want.metricName]))
			}
		})
	}

	// check counter is added to itself
	request := httptest.NewRequest(http.MethodPost, "/update/counter/cnt/10", nil)
	w := httptest.NewRecorder()
	SaveMetric(w, request)
	assert.Equal(t, int64(10+1), container.Instance.MemStorage.Counters["cnt"])

	// check gauge is substituted
	request = httptest.NewRequest(http.MethodPost, "/update/gauge/gaugeName/3", nil)
	w = httptest.NewRecorder()
	SaveMetric(w, request)
	assert.Equal(t, float64(3), container.Instance.MemStorage.Gauges["gaugeName"])
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

func Test_validateMetricName(t *testing.T) {
	type args struct {
		nameRaw string
	}
	possibleError := errors.New("name must be a non empty string")

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "ok",
			args:    args{nameRaw: "name"},
			wantErr: false,
		},
		{
			name:    "empty",
			args:    args{nameRaw: ""},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateMetricName(tt.args.nameRaw); (err != nil) != tt.wantErr {
				t.Errorf("validateMetricName() error = %v, wantErr %v", err, possibleError)
			}
		})
	}
}

func Test_validateMetricType(t *testing.T) {
	type args struct {
		metricTypeRaw string
		statusCode    int
	}

	tests := []struct {
		name         string
		args         args
		wantErr      bool
		errorMessage string
	}{
		{
			name:    "counter",
			args:    args{metricTypeRaw: "counter"},
			wantErr: false,
		},
		{
			name:    "gauge",
			args:    args{metricTypeRaw: "gauge"},
			wantErr: false,
		},
		{
			name:    "unknown",
			args:    args{metricTypeRaw: "unknown"},
			wantErr: true,
		},
		{
			name:    "wrong",
			args:    args{metricTypeRaw: "wrong"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateMetricType(tt.args.metricTypeRaw); (err != nil) != tt.wantErr {
				t.Errorf("validateMetricType() error = %v, wantErr %v", err, tt.wantErr)
				if tt.errorMessage != "" {
					require.Equal(t, tt.errorMessage, err.Error())
				}
			}
		})
	}
}

func Test_validateMetricValue(t *testing.T) {
	type args struct {
		metricTypeValidated string
		valueRaw            string
	}
	tests := []struct {
		name             string
		args             args
		wantCounterValue int64
		wantGaugeValue   float64
		wantErr          bool
	}{
		{
			name:             "counter ok",
			args:             args{metricTypeValidated: "counter", valueRaw: "10"},
			wantCounterValue: 10,
			wantGaugeValue:   0,
			wantErr:          false,
		},
		{
			name:             "gauge ok",
			args:             args{metricTypeValidated: "gauge", valueRaw: "10.45"},
			wantCounterValue: 0,
			wantGaugeValue:   10.45,
			wantErr:          false,
		},

		{
			name:             "counter float error",
			args:             args{metricTypeValidated: "counter", valueRaw: "10.45"},
			wantCounterValue: 10,
			wantGaugeValue:   0,
			wantErr:          true,
		},

		{
			name:             "counter string error",
			args:             args{metricTypeValidated: "counter", valueRaw: "hello"},
			wantCounterValue: 0,
			wantGaugeValue:   0,
			wantErr:          true,
		},
		{
			name:             "gauge string error",
			args:             args{metricTypeValidated: "gauge", valueRaw: "hello"},
			wantCounterValue: 0,
			wantGaugeValue:   0,
			wantErr:          true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCounterValue, gotGaugeValue, err := validateMetricValue(tt.args.metricTypeValidated, tt.args.valueRaw)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("validateMetricValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotCounterValue != tt.wantCounterValue {
				t.Errorf("validateMetricValue() gotCounterValue = %v, want %v", gotCounterValue, tt.wantCounterValue)
			}
			if gotGaugeValue != tt.wantGaugeValue {
				t.Errorf("validateMetricValue() gotGaugeValue = %v, want %v", gotGaugeValue, tt.wantGaugeValue)
			}
		})
	}
}

func Test_validateParameters(t *testing.T) {
	type args struct {
		metricTypeRaw string
		nameRaw       string
		valueRaw      string
	}
	tests := []struct {
		name             string
		args             args
		wantCounterValue int64
		wantGaugeValue   float64
		wantErr          bool
	}{
		{
			name:             "counter",
			args:             args{metricTypeRaw: "counter", nameRaw: "name", valueRaw: "1"},
			wantCounterValue: 1,
			wantGaugeValue:   0,
			wantErr:          false,
		},
		{
			name:             "gauge",
			args:             args{metricTypeRaw: "gauge", nameRaw: "name", valueRaw: "1.1"},
			wantCounterValue: 0,
			wantGaugeValue:   1.1,
			wantErr:          false,
		},
		{
			name:             "gauge string error",
			args:             args{metricTypeRaw: "gauge", nameRaw: "name", valueRaw: "hello"},
			wantCounterValue: 0,
			wantGaugeValue:   0,
			wantErr:          true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCounterValue, gotGaugeValue, err := validateParameters(tt.args.metricTypeRaw, tt.args.nameRaw, tt.args.valueRaw)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateParameters() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotCounterValue != tt.wantCounterValue {
				t.Errorf("validateParameters() gotCounterValue = %v, want %v", gotCounterValue, tt.wantCounterValue)
			}
			if gotGaugeValue != tt.wantGaugeValue {
				t.Errorf("validateParameters() gotGaugeValue = %v, want %v", gotGaugeValue, tt.wantGaugeValue)
			}
		})
	}
}
