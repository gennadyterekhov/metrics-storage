package services

import (
	"context"
	"testing"

	"github.com/gennadyterekhov/metrics-storage/internal/server/config"

	"github.com/gennadyterekhov/metrics-storage/internal/server/services/services"

	"github.com/gennadyterekhov/metrics-storage/internal/common/tests"
	"github.com/stretchr/testify/suite"

	"github.com/gennadyterekhov/metrics-storage/internal/common/constants/types"
	"github.com/gennadyterekhov/metrics-storage/internal/server/httpui/requests"
	"github.com/stretchr/testify/assert"
)

type saveMetricTestSuite struct {
	tests.BaseSuite
	Service services.SaveMetricService
	Config  config.ServerConfig
}

func (suite *saveMetricTestSuite) SetupSuite() {
	tests.InitBaseSuite(suite)
	suite.Config = config.New()
	suite.Service = services.NewSaveMetricService(suite.Repository, &suite.Config)
}

func TestSaveMetricService(t *testing.T) {
	suite.Run(t, new(saveMetricTestSuite))
}

func (st *saveMetricTestSuite) TestSaveMetricToMemory() {
	type args struct {
		metricType   string
		name         string
		counterValue int64
		gaugeValue   float64
	}
	cases := []struct {
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
	for _, tt := range cases {
		st.T().Run(tt.name, func(t *testing.T) {
			filledDto := &requests.SaveMetricRequest{
				MetricType:   tt.args.metricType,
				MetricName:   tt.args.name,
				CounterValue: &tt.args.counterValue,
				GaugeValue:   &tt.args.gaugeValue,
			}
			st.Service.SaveMetricToMemory(context.Background(), filledDto)

			if tt.args.metricType == types.Counter {
				assert.Equal(t, tt.args.counterValue, st.Repository.GetCounterOrZero(context.Background(), tt.args.name))
			}
			if tt.args.metricType == types.Gauge {
				assert.Equal(t, tt.args.gaugeValue, st.Repository.GetGaugeOrZero(context.Background(), tt.args.name))
			}
		})
	}

	// check counter is added to itself
	ten := int64(10)
	zeroInt := int64(0)
	zeroFloat := float64(0.0)
	two := float64(2.5)
	st.Service.SaveMetricToMemory(context.Background(), &requests.SaveMetricRequest{
		MetricType: types.Counter, MetricName: "cnt", CounterValue: &ten, GaugeValue: &zeroFloat,
	})
	assert.Equal(st.T(), int64(10+1), st.Repository.GetCounterOrZero(context.Background(), "cnt"))

	// check gauge is substituted, (not 2.5+1.6)
	st.Service.SaveMetricToMemory(context.Background(), &requests.SaveMetricRequest{
		MetricType: types.Gauge, MetricName: "gaugeName", CounterValue: &zeroInt, GaugeValue: &two,
	})
	assert.Equal(st.T(), 2.5, st.Repository.GetGaugeOrZero(context.Background(), "gaugeName"))
}
