package storage

import (
	"context"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSaveLoad(t *testing.T) {
	ctx := context.Background()
	MetricsRepository.AddCounter(ctx, "c1", 1)
	MetricsRepository.AddCounter(ctx, "c2", 2)
	MetricsRepository.AddCounter(ctx, "c3", 3)
	MetricsRepository.SetGauge(ctx, "g1", 1.5)
	MetricsRepository.SetGauge(ctx, "g2", 2.6)
	MetricsRepository.SetGauge(ctx, "g3", 3.7)

	filename := `metrics.json`
	var err error
	err = MetricsRepository.SaveToDisk(ctx, filename)
	assert.NoError(t, err)
	assert.NoError(t, err)

	var result MemStorage
	err = (&result).LoadFromDisk(ctx, filename)
	assert.NoError(t, err)

	if MetricsRepository.GetDB() == nil {
		assert.True(t, MetricsRepository.GetMemStorage().IsEqual(&result))
	} else {
		assert.True(t, MetricsRepository.GetDB().IsEqual(&result))
	}

	err = os.Remove(filename)
	assert.NoError(t, err)
}

func TestLoadEmptyWhenError(t *testing.T) {
	ctx := context.Background()

	originalRepository := MetricsRepository
	originalRepository.AddCounter(ctx, "c1", 1)

	filename := `metrics.json`
	var err error
	var result MemStorage
	err = (&result).LoadFromDisk(ctx, filename)
	assert.Error(t, err)

	_, err = result.GetCounter(ctx, "c1")
	assert.Error(t, err)

	if MetricsRepository.GetDB() == nil {
		assert.True(t, MetricsRepository.GetMemStorage().IsEqual(&result))
	} else {
		assert.True(t, MetricsRepository.GetDB().IsEqual(&result))
	}
}

func (strg *DBStorage) IsEqual(anotherStorage Interface) (eq bool) {
	ctx := context.Background()

	gauges, counters := strg.GetAllGauges(ctx), strg.GetAllCounters(ctx)
	gauges2, counters2 := anotherStorage.GetAllGauges(ctx), anotherStorage.GetAllCounters(ctx)

	return reflect.DeepEqual(gauges, gauges2) && reflect.DeepEqual(counters, counters2)
}

func (strg *MemStorage) IsEqual(anotherStorage Interface) (eq bool) {
	ctx := context.Background()

	gauges, counters := strg.GetAllGauges(ctx), strg.GetAllCounters(ctx)
	gauges2, counters2 := anotherStorage.GetAllGauges(ctx), anotherStorage.GetAllCounters(ctx)

	return reflect.DeepEqual(gauges, gauges2) && reflect.DeepEqual(counters, counters2)
}
