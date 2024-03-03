package agent

import (
	"github.com/gennadyterekhov/metrics-storage/internal/server/httpui/handlers"
	"github.com/gennadyterekhov/metrics-storage/internal/server/storage"
	"github.com/gennadyterekhov/metrics-storage/internal/testhelper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestMain(m *testing.M) {
	testhelper.BootstrapWithServer(
		m,
		httptest.NewServer(
			handlers.GetRouter(),
		),
	)
}

func TestAgent(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "test",
		},
	}

	oneIteration := 1
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := RunAgent(&AgentConfig{
				Addr:            testhelper.TestServer.URL,
				ReportInterval:  1,
				PollInterval:    1,
				TotalIterations: &oneIteration,
			})
			require.NoError(t, err)

			assert.Equal(t,
				1,
				len(storage.MetricsRepository.GetAllCounters()),
			)
			assert.Equal(t,
				27+1,
				len(storage.MetricsRepository.GetAllGauges()),
			)
		})
	}
}

func TestGzip(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "test",
		},
	}

	oneIteration := 1
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := RunAgent(&AgentConfig{
				Addr:            testhelper.TestServer.URL,
				ReportInterval:  1,
				PollInterval:    1,
				IsGzip:          true,
				TotalIterations: &oneIteration,
			})
			require.NoError(t, err)

			assert.Equal(t,
				1,
				len(storage.MetricsRepository.GetAllCounters()),
			)
			assert.Equal(t,
				27+1,
				len(storage.MetricsRepository.GetAllGauges()),
			)
			savedValue := storage.MetricsRepository.GetCounterOrZero("PollCount")
			assert.Equal(t, int64(2), savedValue)
		})
	}
}

func TestSameValueReturnedFromServer(t *testing.T) {

	tests := []struct {
		name string
	}{
		{
			name: "test",
		},
	}
	oneIteration := 1

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := RunAgent(&AgentConfig{
				Addr:            testhelper.TestServer.URL,
				ReportInterval:  1,
				PollInterval:    1,
				TotalIterations: &oneIteration,
			})
			require.NoError(t, err)

			assert.Equal(t,
				1,
				len(storage.MetricsRepository.GetAllCounters()),
			)
			assert.Equal(t,
				27+1,
				len(storage.MetricsRepository.GetAllGauges()),
			)

			url := "/value/gauge/BuckHashSys"

			_, responseBody := testhelper.SendRequest(
				t,
				testhelper.TestServer,
				http.MethodGet,
				url,
			)

			savedValue := storage.MetricsRepository.GetGaugeOrZero("BuckHashSys")

			assert.Equal(
				t,
				strconv.FormatFloat(savedValue, 'g', -1, 64),
				string(responseBody),
			)

		})
	}
}
