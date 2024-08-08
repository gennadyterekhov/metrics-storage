package storage

import (
	"context"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestSaveLoad(t *testing.T) {
	ctx := context.Background()
	repo := New("")
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

func TestSaveLoadDBUsingContainer(t *testing.T) {
	ctx := context.Background()
	cont, err := createPsqlContainer(ctx)
	assert.NoError(t, err)

	repo := New(TestDBDsn)
	assertCanReadAndWriteToDisk(ctx, t, repo)

	assert.NoError(t, cont.Terminate(ctx))
}

func isEqual(strg Interface, anotherStorage Interface) (eq bool) {
	ctx := context.Background()

	gauges, counters := strg.GetAllGauges(ctx), strg.GetAllCounters(ctx)
	gauges2, counters2 := anotherStorage.GetAllGauges(ctx), anotherStorage.GetAllCounters(ctx)

	return reflect.DeepEqual(gauges, gauges2) && reflect.DeepEqual(counters, counters2)
}

func assertCanReadAndWriteToDisk(ctx context.Context, t *testing.T, repo Interface) {
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

func createPsqlContainer(ctx context.Context) (*postgres.PostgresContainer, error) {
	dbName := "metrics_db_test"
	dbUser := "metrics_user"
	dbPassword := "metrics_pass"

	postgresContainer, err := postgres.Run(ctx,
		"docker.io/postgres:16-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)

	return postgresContainer, err
}
