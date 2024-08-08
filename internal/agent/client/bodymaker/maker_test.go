package bodymaker

import (
	"testing"

	"github.com/gennadyterekhov/metrics-storage/internal/agent/metric"
	"github.com/gennadyterekhov/metrics-storage/internal/common/constants/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncryption(t *testing.T) {
	publicKeyFilePath := "../../../../keys/public.test"

	metrics := metric.CounterMetric{
		Name:  "nm",
		Type:  types.Counter,
		Value: 1,
	}
	metricsJSON := `{"id":"nm","type":"counter","delta":1,"value":0}`
	bodyBytes, err := GetBody(&metrics, publicKeyFilePath)
	require.NoError(t, err)

	assert.NotEqual(t,
		metricsJSON,
		string(bodyBytes),
	)

	bodyBytes, err = GetBody(&metrics, "")
	require.NoError(t, err)

	assert.Equal(t,
		metricsJSON,
		string(bodyBytes),
	)
}
