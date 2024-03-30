package storage

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestSaveLoad(t *testing.T) {
	MetricsRepository.AddCounter("c1", 1)
	MetricsRepository.AddCounter("c2", 2)
	MetricsRepository.AddCounter("c3", 3)
	MetricsRepository.SetGauge("g1", 1.5)
	MetricsRepository.SetGauge("g2", 2.6)
	MetricsRepository.SetGauge("g3", 3.7)

	filename := `metrics.json`
	var err error
	err = MetricsRepository.SaveToDisk(filename)
	assert.NoError(t, err)
	assert.NoError(t, err)

	var result MemStorage
	err = (&result).LoadFromDisk(filename)
	assert.NoError(t, err)

	assert.True(t, MetricsRepository.IsEqual(&result))

	err = os.Remove(filename)
	assert.NoError(t, err)

}

func TestLoadEmptyWhenError(t *testing.T) {
	originalRepository := MetricsRepository
	originalRepository.AddCounter("c1", 1)

	filename := `metrics.json`
	var err error
	var result MemStorage
	err = (&result).LoadFromDisk(filename)
	assert.Error(t, err)

	_, err = result.GetCounter("c1")
	assert.Error(t, err)

	assert.False(t, originalRepository.IsEqual(&result))
}
