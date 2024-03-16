package agent

import (
	"context"
	"github.com/gennadyterekhov/metrics-storage/internal/server/httpui/handlers"
	"github.com/gennadyterekhov/metrics-storage/internal/server/storage"
	"github.com/gennadyterekhov/metrics-storage/internal/testhelper"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strconv"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	testhelper.BootstrapWithDefaultServer(m, handlers.GetRouter())
}

func TestAgent(t *testing.T) {
	storage.MetricsRepository.Clear()

	tests := []struct {
		name string
	}{
		{
			name: "test",
		},
	}
	ctx, cancelContextFn := context.WithTimeout(context.Background(), 100*time.Millisecond)

	defer cancelContextFn()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go runAgentRoutine(ctx, &AgentConfig{
				Addr:           testhelper.TestServer.URL,
				ReportInterval: 1,
				PollInterval:   1,
			})

			<-ctx.Done()

			contextEndCondition := ctx.Err()

			if contextEndCondition == context.DeadlineExceeded || contextEndCondition == context.Canceled {
				assert.Equal(t,
					1,
					len(storage.MetricsRepository.GetAllCounters()),
				)
				assert.Equal(t,
					27+1,
					len(storage.MetricsRepository.GetAllGauges()),
				)

				return
			} else {
				t.Error("context didnt finish")
			}

		})
	}
}

func TestBatch(t *testing.T) {
	storage.MetricsRepository.Clear()

	ctx, cancelContextFn := context.WithTimeout(context.Background(), 100*time.Millisecond)

	defer cancelContextFn()

	t.Run("batch", func(t *testing.T) {
		go runAgentRoutine(ctx, &AgentConfig{
			Addr:           testhelper.TestServer.URL,
			ReportInterval: 1,
			PollInterval:   1,
			IsBatch:        true,
		})

		<-ctx.Done()

		contextEndCondition := ctx.Err()

		if contextEndCondition == context.DeadlineExceeded || contextEndCondition == context.Canceled {
			assert.Equal(t,
				1,
				len(storage.MetricsRepository.GetAllCounters()),
			)
			assert.Equal(t,
				27+1,
				len(storage.MetricsRepository.GetAllGauges()),
			)

			return
		} else {
			t.Error("context didnt finish")
		}
	})

}

func TestGzip(t *testing.T) {
	storage.MetricsRepository.Clear()

	tests := []struct {
		name string
	}{
		{
			name: "test",
		},
	}
	ctx, cancelContextFn := context.WithTimeout(context.Background(), 100*time.Millisecond)

	defer cancelContextFn()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go runAgentRoutine(ctx, &AgentConfig{
				Addr:           testhelper.TestServer.URL,
				ReportInterval: 1,
				PollInterval:   1,
				IsGzip:         true,
			})

			<-ctx.Done()

			contextEndCondition := ctx.Err()

			if contextEndCondition == context.DeadlineExceeded || contextEndCondition == context.Canceled {
				assert.Equal(t,
					1,
					len(storage.MetricsRepository.GetAllCounters()),
				)
				assert.Equal(t,
					27+1,
					len(storage.MetricsRepository.GetAllGauges()),
				)
				savedValue := storage.MetricsRepository.GetCounterOrZero("PollCount")
				assert.Equal(t, int64(1), savedValue)
				return
			}

			t.Error("context didnt finish")
		})
	}
}

func TestSameValueReturnedFromServer(t *testing.T) {
	storage.MetricsRepository.Clear()

	tests := []struct {
		name string
	}{
		{
			name: "test",
		},
	}
	ctx, cancelContextFn := context.WithTimeout(context.Background(), 100*time.Millisecond)

	defer cancelContextFn()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go runAgentRoutine(ctx, &AgentConfig{
				Addr:           testhelper.TestServer.URL,
				ReportInterval: 1,
				PollInterval:   1,
			})

			<-ctx.Done()

			contextEndCondition := ctx.Err()

			if contextEndCondition == context.DeadlineExceeded || contextEndCondition == context.Canceled {
				assert.Equal(t,
					1,
					len(storage.MetricsRepository.GetAllCounters()),
				)
				assert.Equal(t,
					27+1,
					len(storage.MetricsRepository.GetAllGauges()),
				)

				url := "/value/gauge/BuckHashSys"

				r, responseBody := testhelper.SendRequest(
					t,
					testhelper.TestServer,
					http.MethodGet,
					url,
				)
				r.Body.Close()
				savedValue := storage.MetricsRepository.GetGaugeOrZero("BuckHashSys")

				assert.Equal(
					t,
					strconv.FormatFloat(savedValue, 'g', -1, 64),
					string(responseBody),
				)

				return
			} else {

				t.Error("context didnt finish")
			}

		})
	}
}

func runAgentRoutine(ctx context.Context, config *AgentConfig) {
	RunAgent(config)
}
