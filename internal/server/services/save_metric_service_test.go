package services

import (
	"context"
	"testing"

	"github.com/gennadyterekhov/metrics-storage/internal/common/constants/types"
	"github.com/gennadyterekhov/metrics-storage/internal/server/httpui/requests"
	"github.com/gennadyterekhov/metrics-storage/internal/server/storage"
	"github.com/stretchr/testify/assert"
)

func TestSaveMetricToMemory(t *testing.T) {
	type args struct {
		metricType   string
		name         string
		counterValue int64
		gaugeValue   float64
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Counter",
			args: args{metricType: "counter", name: "cnt", counterValue: int64(1), gaugeValue: 0},
		},
		{
			name: "Gauge",
			args: args{metricType: "gauge", name: "gaugeName", counterValue: 0, gaugeValue: 1.6},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filledDto := &requests.SaveMetricRequest{
				MetricType:   tt.args.metricType,
				MetricName:   tt.args.name,
				CounterValue: &tt.args.counterValue,
				GaugeValue:   &tt.args.gaugeValue,
			}
			SaveMetricToMemory(context.Background(), filledDto)

			if tt.args.metricType == types.Counter {
				assert.Equal(t, tt.args.counterValue, storage.MetricsRepository.GetCounterOrZero(context.Background(), tt.args.name))
			}
			if tt.args.metricType == types.Gauge {
				assert.Equal(t, tt.args.gaugeValue, storage.MetricsRepository.GetGaugeOrZero(context.Background(), tt.args.name))
			}
		})
	}

	// check counter is added to itself
	ten := int64(10)
	zeroInt := int64(0)
	zeroFloat := float64(0.0)
	two := float64(2.5)
	SaveMetricToMemory(context.Background(), &requests.SaveMetricRequest{
		MetricType: types.Counter, MetricName: "cnt", CounterValue: &ten, GaugeValue: &zeroFloat,
	})
	assert.Equal(t, int64(10+1), storage.MetricsRepository.GetCounterOrZero(context.Background(), "cnt"))

	// check gauge is substituted, (not 2.5+1.6)
	SaveMetricToMemory(context.Background(), &requests.SaveMetricRequest{
		MetricType: types.Gauge, MetricName: "gaugeName", CounterValue: &zeroInt, GaugeValue: &two,
	})
	assert.Equal(t, 2.5, storage.MetricsRepository.GetGaugeOrZero(context.Background(), "gaugeName"))
}
