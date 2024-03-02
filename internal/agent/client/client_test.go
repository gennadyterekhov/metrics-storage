package client

import (
	"github.com/gennadyterekhov/metrics-storage/internal/agent/metric"
	"github.com/gennadyterekhov/metrics-storage/internal/constants/types"
	"github.com/gennadyterekhov/metrics-storage/internal/container"
	"github.com/gennadyterekhov/metrics-storage/internal/server/httpui/handlers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http/httptest"
	"testing"
)

func TestCanSendCounterValue(t *testing.T) {
	testServer := httptest.NewServer(
		handlers.GetRouter(),
	)

	metricsStorageClient := MetricsStorageClient{
		Address:          testServer.URL,
		IsGzip:           false,
		SendWhenNoServer: false,
	}

	type want struct {
		counterValue int64
	}
	tests := []struct {
		name   string
		isGzip bool
		want   want
	}{
		{
			name:   "send one counter",
			isGzip: false,
			want:   want{10},
		},
		{
			name:   "send one counter gzip",
			isGzip: true,
			want:   want{10},
		},
	}
	var err error
	for _, tt := range tests {
		container.MetricsRepository.Clear()
		t.Run(tt.name, func(t *testing.T) {
			if tt.isGzip {
				metricsStorageClient.IsGzip = true
			}
			metrics := metric.CounterMetric{
				Name:  "nm",
				Type:  types.Counter,
				Value: tt.want.counterValue,
			}
			err = metricsStorageClient.SendMetric(&metrics)
			require.NoError(t, err)

			assert.Equal(t,
				1,
				len(container.MetricsRepository.GetAllCounters()),
			)
			assert.Equal(t,
				0,
				len(container.MetricsRepository.GetAllGauges()),
			)

			assert.Equal(t,
				tt.want.counterValue,
				container.MetricsRepository.GetCounterOrZero("nm"),
			)
		})
	}
}

func TestCanSendGaugeValue(t *testing.T) {
	container.MetricsRepository.Clear()
	testServer := httptest.NewServer(
		handlers.GetRouter(),
	)

	metricsStorageClient := MetricsStorageClient{
		Address:          testServer.URL,
		IsGzip:           false,
		SendWhenNoServer: false,
	}
	type want struct {
		gaugeValue float64
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "send one gauge",
			want: want{5.5},
		},
	}
	var err error
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metrics := metric.GaugeMetric{
				Name:  "nm",
				Type:  types.Gauge,
				Value: tt.want.gaugeValue,
			}
			err = metricsStorageClient.SendMetric(&metrics)

			require.NoError(t, err)

			assert.Equal(t,
				0,
				len(container.MetricsRepository.GetAllCounters()),
			)
			assert.Equal(t,
				1,
				len(container.MetricsRepository.GetAllGauges()),
			)

			assert.Equal(t,
				tt.want.gaugeValue,
				container.MetricsRepository.GetGaugeOrZero("nm"),
			)
		})
	}
}
