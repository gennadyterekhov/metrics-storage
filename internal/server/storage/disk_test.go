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
	repo := NewRAMStorage()
	repo.AddCounter(ctx, "c1", 1)
	repo.AddCounter(ctx, "c2", 2)
	repo.AddCounter(ctx, "c3", 3)
	repo.SetGauge(ctx, "g1", 1.5)
	repo.SetGauge(ctx, "g2", 2.6)
	repo.SetGauge(ctx, "g3", 3.7)

	filename := `metrics.json`
	var err error
	err = repo.SaveToDisk(ctx, filename)
	assert.NoError(t, err)
	assert.NoError(t, err)

	var result MemStorage
	err = (&result).LoadFromDisk(ctx, filename)
	assert.NoError(t, err)

	assert.True(t, isEqual(repo, &result))

	err = os.Remove(filename)
	assert.NoError(t, err)
}

func TestLoadEmptyWhenError(t *testing.T) {
	ctx := context.Background()

	repo := NewRAMStorage()
	repo.AddCounter(ctx, "c1", 1)

	filename := `metrics.json`
	var err error
	var result MemStorage
	err = (&result).LoadFromDisk(ctx, filename)
	assert.Error(t, err)

	_, err = result.GetCounter(ctx, "c1")
	assert.Error(t, err)

	assert.False(t, isEqual(repo, &result))
	assert.True(t, isEqual(NewRAMStorage(), &result))
}

func isEqual(strg *MemStorage, anotherStorage StorageInterface) (eq bool) {
	ctx := context.Background()

	gauges, counters := strg.GetAllGauges(ctx), strg.GetAllCounters(ctx)
	gauges2, counters2 := anotherStorage.GetAllGauges(ctx), anotherStorage.GetAllCounters(ctx)

	return reflect.DeepEqual(gauges, gauges2) && reflect.DeepEqual(counters, counters2)
}
