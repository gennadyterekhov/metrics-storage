package services

import (
	"github.com/gennadyterekhov/metrics-storage/internal/container"
	"github.com/gennadyterekhov/metrics-storage/internal/dto"
	"github.com/gennadyterekhov/metrics-storage/internal/types"
	"github.com/stretchr/testify/assert"
	"testing"
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
			filledDto := &dto.MetricToSaveDto{
				Type:         tt.args.metricType,
				Name:         tt.args.name,
				CounterValue: tt.args.counterValue,
				GaugeValue:   tt.args.gaugeValue,
			}
			SaveMetricToMemory(filledDto)

			if tt.args.metricType == types.Counter {
				assert.Equal(t, tt.args.counterValue, container.MetricsRepository.GetCounterOrZero(tt.args.name))
			}
			if tt.args.metricType == types.Gauge {
				assert.Equal(t, tt.args.gaugeValue, container.MetricsRepository.GetGaugeOrZero(tt.args.name))
			}
		})
	}

	// check counter is added to itself
	SaveMetricToMemory(&dto.MetricToSaveDto{
		Type: types.Counter, Name: "cnt", CounterValue: 10, GaugeValue: 0,
	})
	assert.Equal(t, int64(10+1), container.MetricsRepository.GetCounterOrZero("cnt"))

	// check gauge is substituted, (not 2.5+1.6)
	SaveMetricToMemory(&dto.MetricToSaveDto{
		Type: types.Gauge, Name: "gaugeName", CounterValue: 0, GaugeValue: 2.5,
	})
	assert.Equal(t, 2.5, container.MetricsRepository.GetGaugeOrZero("gaugeName"))
}
