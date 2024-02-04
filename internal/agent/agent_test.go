package agent

import (
	"github.com/gennadyterekhov/metrics-storage/internal/container"
	"github.com/gennadyterekhov/metrics-storage/internal/router"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http/httptest"
	"testing"
)

func shouldContinueMock(iter int) bool {
	return iter == 0
}

func TestAgent(t *testing.T) {
	testServer := httptest.NewServer(
		router.GetRouter(),
	)

	url := testServer.URL
	tests := []struct {
		name string
	}{
		{
			name: "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Agent(url, shouldContinueMock)
			require.NoError(t, err)

			assert.Equal(t,
				1,
				len(container.Instance.MetricsRepository.GetAllCounters()),
			)
			assert.Equal(t,
				27+1,
				len(container.Instance.MetricsRepository.GetAllGauges()),
			)
		})
	}
}
