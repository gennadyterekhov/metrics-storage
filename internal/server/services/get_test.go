package services

import (
	"context"
	"fmt"
	"testing"

	"github.com/gennadyterekhov/metrics-storage/internal/server/repositories"
	"github.com/gennadyterekhov/metrics-storage/internal/server/storage"

	"github.com/gennadyterekhov/metrics-storage/internal/common/constants/types"
	"github.com/gennadyterekhov/metrics-storage/internal/common/tests"
	"github.com/gennadyterekhov/metrics-storage/internal/server/http/requests"
	"github.com/gennadyterekhov/metrics-storage/internal/server/services/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type getMetricServiceSuite struct {
	tests.BaseSuite
	Service *services.GetMetricService
}

func (suite *getMetricServiceSuite) SetupSuite() {
	tests.InitBaseSuite(suite)
	suite.Service = services.NewGetMetricService(suite.GetRepository())
}

func BenchmarkGetMetricService(b *testing.B) {
	repo := repositories.New(storage.New(""))
	srv := services.NewGetMetricService(repo)
	for i := 0; i < b.N; i++ {
		repo.SetGauge(context.Background(), fmt.Sprintf("g%d", i), float64(i))
		repo.AddCounter(context.Background(), fmt.Sprintf("c%d", i), int64(i))
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		requestDto := requests.GetMetricRequest{MetricType: types.Gauge, MetricName: fmt.Sprintf("g%d", b.N), IsJSON: true, Error: nil}
		srv.GetMetric(context.Background(), &requestDto)

		requestDto = requests.GetMetricRequest{MetricType: types.Counter, MetricName: fmt.Sprintf("c%d", b.N), IsJSON: true, Error: nil}
		srv.GetMetric(context.Background(), &requestDto)
	}
}

func TestGetMetricService(t *testing.T) {
	suite.Run(t, new(getMetricServiceSuite))
}

func (suite *getMetricServiceSuite) TestCanGetGaugeAndCounter() {
	suite.Repository.SetGauge(context.Background(), "g1", float64(1))
	suite.Repository.SetGauge(context.Background(), "g2", float64(2))
	suite.Repository.SetGauge(context.Background(), "g3", float64(3))

	requestDto := requests.GetMetricRequest{MetricType: types.Gauge, MetricName: "g1", IsJSON: true, Error: nil}
	responseDto := suite.Service.GetMetric(context.Background(), &requestDto)

	assert.NoError(suite.T(), responseDto.Error)
	assert.Equal(suite.T(), float64(1), *responseDto.GaugeValue)

	requestDto = requests.GetMetricRequest{MetricType: types.Gauge, MetricName: "g2", IsJSON: true, Error: nil}
	responseDto = suite.Service.GetMetric(context.Background(), &requestDto)

	assert.NoError(suite.T(), responseDto.Error)
	assert.Equal(suite.T(), float64(2), *responseDto.GaugeValue)

	requestDto = requests.GetMetricRequest{MetricType: types.Gauge, MetricName: "g3", IsJSON: true, Error: nil}
	responseDto = suite.Service.GetMetric(context.Background(), &requestDto)

	assert.NoError(suite.T(), responseDto.Error)
	assert.Equal(suite.T(), float64(3), *responseDto.GaugeValue)
}

func (suite *getMetricServiceSuite) TestCanGetCounter() {
	suite.Repository.AddCounter(context.Background(), "c1", int64(1))
	suite.Repository.AddCounter(context.Background(), "c2", int64(2))
	suite.Repository.AddCounter(context.Background(), "c3", int64(3))

	requestDto := requests.GetMetricRequest{MetricType: types.Counter, MetricName: "c1", IsJSON: true, Error: nil}
	responseDto := suite.Service.GetMetric(context.Background(), &requestDto)

	assert.NoError(suite.T(), responseDto.Error)
	assert.Equal(suite.T(), int64(1), *responseDto.CounterValue)

	requestDto = requests.GetMetricRequest{MetricType: types.Counter, MetricName: "c2", IsJSON: true, Error: nil}
	responseDto = suite.Service.GetMetric(context.Background(), &requestDto)

	assert.NoError(suite.T(), responseDto.Error)
	assert.Equal(suite.T(), int64(2), *responseDto.CounterValue)

	requestDto = requests.GetMetricRequest{MetricType: types.Counter, MetricName: "c3", IsJSON: true, Error: nil}
	responseDto = suite.Service.GetMetric(context.Background(), &requestDto)

	assert.NoError(suite.T(), responseDto.Error)
	assert.Equal(suite.T(), int64(3), *responseDto.CounterValue)
}
