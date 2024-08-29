package services

import (
	"context"
	"fmt"
	"math/rand"
	"testing"

	"github.com/gennadyterekhov/metrics-storage/internal/server/repositories"
	"github.com/gennadyterekhov/metrics-storage/internal/server/storage"

	"github.com/gennadyterekhov/metrics-storage/internal/server/config"

	"github.com/gennadyterekhov/metrics-storage/internal/server/services/services"

	"github.com/gennadyterekhov/metrics-storage/internal/common/tests"
	"github.com/stretchr/testify/suite"

	"github.com/gennadyterekhov/metrics-storage/internal/common/constants/types"
	"github.com/gennadyterekhov/metrics-storage/internal/server/http/requests"
	"github.com/stretchr/testify/assert"
)

type saveMetricTestSuite struct {
	tests.BaseSuite
	Service *services.SaveMetricService
	Config  *config.ServerConfig
}

func (suite *saveMetricTestSuite) SetupSuite() {
	tests.InitBaseSuite(suite)
	suite.Config = config.New()
	suite.Service = services.NewSaveMetricService(suite.Repository, suite.Config)
}

func BenchmarkSaveMetricService(b *testing.B) {
	conf := config.New()
	repo := repositories.New(storage.New(""))
	srv := services.NewSaveMetricService(repo, conf)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rndVal1 := rand.Int63()
		rndVal2 := rand.Float64()
		srv.SaveMetricToMemory(context.Background(), &requests.SaveMetricRequest{
			MetricType: types.Counter, MetricName: fmt.Sprintf("c%d", i), CounterValue: &rndVal1, GaugeValue: nil,
		})

		srv.SaveMetricToMemory(context.Background(), &requests.SaveMetricRequest{
			MetricType: types.Gauge, MetricName: fmt.Sprintf("g%d", i), CounterValue: nil, GaugeValue: &rndVal2,
		})
	}
}

func TestSaveMetricService(t *testing.T) {
	suite.Run(t, new(saveMetricTestSuite))
}

func (suite *saveMetricTestSuite) TestSaveMetricToMemory() {
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
		suite.T().Run(tt.name, func(t *testing.T) {
			filledDto := &requests.SaveMetricRequest{
				MetricType:   tt.args.metricType,
				MetricName:   tt.args.name,
				CounterValue: &tt.args.counterValue,
				GaugeValue:   &tt.args.gaugeValue,
			}
			suite.Service.SaveMetricToMemory(context.Background(), filledDto)

			if tt.args.metricType == types.Counter {
				assert.Equal(t, tt.args.counterValue, suite.Repository.GetCounterOrZero(context.Background(), tt.args.name))
			}
			if tt.args.metricType == types.Gauge {
				assert.Equal(t, tt.args.gaugeValue, suite.Repository.GetGaugeOrZero(context.Background(), tt.args.name))
			}
		})
	}

	// check counter is added to itself
	ten := int64(10)
	two := float64(2.5)
	suite.Service.SaveMetricToMemory(context.Background(), &requests.SaveMetricRequest{
		MetricType: types.Counter, MetricName: "cnt", CounterValue: &ten, GaugeValue: nil,
	})
	assert.Equal(suite.T(), int64(10+1), suite.Repository.GetCounterOrZero(context.Background(), "cnt"))

	// check gauge is substituted, (not 2.5+1.6)
	suite.Service.SaveMetricToMemory(context.Background(), &requests.SaveMetricRequest{
		MetricType: types.Gauge, MetricName: "gaugeName", CounterValue: nil, GaugeValue: &two,
	})
	assert.Equal(suite.T(), 2.5, suite.Repository.GetGaugeOrZero(context.Background(), "gaugeName"))
}
