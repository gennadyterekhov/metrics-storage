package agent

import (
	"context"
	"github.com/gennadyterekhov/metrics-storage/internal/container"
	"github.com/gennadyterekhov/metrics-storage/internal/handlers"
	"github.com/gennadyterekhov/metrics-storage/internal/types"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestAgent(t *testing.T) {
	t.Skipf("cannot do it yet")
	testServer := httptest.NewServer(
		handlers.GetRouter(),
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
			err := RunAgent(&AgentConfig{
				Addr:           url,
				ReportInterval: 1,
				PollInterval:   1,
			})
			require.NoError(t, err)

			assert.Equal(t,
				1,
				len(container.MetricsRepository.GetAllCounters()),
			)
			assert.Equal(t,
				27+1,
				len(container.MetricsRepository.GetAllGauges()),
			)
		})
	}
}

func TestSameValueReturnedFromServer(t *testing.T) {
	t.Skipf("cannot do it yet")

	testServer := httptest.NewServer(
		handlers.GetRouter(),
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
			err := RunAgent(&AgentConfig{
				Addr:           url,
				ReportInterval: 1,
				PollInterval:   1,
			})
			require.NoError(t, err)

			assert.Equal(t,
				1,
				len(container.MetricsRepository.GetAllCounters()),
			)
			assert.Equal(t,
				27+1,
				len(container.MetricsRepository.GetAllGauges()),
			)

			url := "/value/gauge/BuckHashSys"
			request := httptest.NewRequest(http.MethodGet, url, nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("metricType", types.Gauge)
			rctx.URLParams.Add("metricName", "BuckHashSys")
			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))
			w := httptest.NewRecorder()
			handlers.GetMetric(w, request)

			res := w.Result()
			metricFromResponse, _ := io.ReadAll(res.Body)
			savedValue := container.MetricsRepository.GetGaugeOrZero("BuckHashSys")

			defer res.Body.Close()
			assert.Equal(
				t,
				strconv.FormatFloat(savedValue, 'g', -1, 64),
				string(metricFromResponse),
			)

		})
	}
}
