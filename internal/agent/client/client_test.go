package client

import (
	"context"
	"testing"

	"github.com/gennadyterekhov/metrics-storage/internal/common/tests"
	"github.com/stretchr/testify/suite"

	"github.com/gennadyterekhov/metrics-storage/internal/agent/metric"
	"github.com/gennadyterekhov/metrics-storage/internal/common/constants/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type clientTestSuite struct {
	tests.BaseSuiteWithServer
}

func (suite *clientTestSuite) SetupSuite() {
	tests.InitBaseSuiteWithServer(suite)
}

func TestAgentSuite(t *testing.T) {
	suite.Run(t, new(clientTestSuite))
}

func (suite *clientTestSuite) TestCanSendCounterValue() {
	metricsStorageClient := MetricsStorageClient{
		Address: suite.TestHTTPServer.Server.URL,
		IsGzip:  false,
	}

	type want struct {
		counterValue int64
	}
	cases := []struct {
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
	for _, tt := range cases {
		suite.T().Run(tt.name, func(t *testing.T) {
			suite.Repository.Clear()
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
				len(suite.Repository.GetAllCounters(context.Background())),
			)
			assert.Equal(t,
				0,
				len(suite.Repository.GetAllGauges(context.Background())),
			)

			assert.Equal(t,
				tt.want.counterValue,
				suite.Repository.GetCounterOrZero(context.Background(), "nm"),
			)
		})
	}
}

func (suite *clientTestSuite) TestCanSendGaugeValue() {
	metricsStorageClient := MetricsStorageClient{
		Address: suite.TestHTTPServer.Server.URL,
		IsGzip:  false,
	}
	type want struct {
		gaugeValue float64
	}
	cases := []struct {
		name string
		want want
	}{
		{
			name: "send one gauge",
			want: want{5.5},
		},
	}
	var err error
	for _, tt := range cases {
		suite.T().Run(tt.name, func(t *testing.T) {
			metrics := metric.GaugeMetric{
				Name:  "nm",
				Type:  types.Gauge,
				Value: tt.want.gaugeValue,
			}
			err = metricsStorageClient.SendMetric(&metrics)

			require.NoError(t, err)

			assert.Equal(t,
				0,
				len(suite.Repository.GetAllCounters(context.Background())),
			)
			assert.Equal(t,
				1,
				len(suite.Repository.GetAllGauges(context.Background())),
			)

			assert.Equal(t,
				tt.want.gaugeValue,
				suite.Repository.GetGaugeOrZero(context.Background(), "nm"),
			)
		})
	}
}
