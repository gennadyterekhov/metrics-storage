package agent

import (
	"github.com/gennadyterekhov/metrics-storage/internal/container"
	"github.com/gennadyterekhov/metrics-storage/internal/handlers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAgent(t *testing.T) {
	testServer := httptest.NewServer(
		http.HandlerFunc(handlers.SaveMetric),
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
			err := Agent(url, func() bool { return false })
			require.NoError(t, err)

			assert.Equal(t,
				1,
				len(container.Instance.MetricsRepository.GetAllCounters()),
			)
			assert.Equal(t,
				14,
				len(container.Instance.MetricsRepository.GetAllGauges()),
			)
		})
	}
}
