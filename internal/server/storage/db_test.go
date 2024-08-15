package storage

import (
	"context"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
)

const TestDBDsn = "host=localhost user=metrics_user password=metrics_pass dbname=metrics_db_test sslmode=disable"

func TestSaveLoadDBUsingContainer(t *testing.T) {
	ctx := context.Background()
	cont, err := createPsqlContainer(ctx)
	assert.NoError(t, err)

	repo := New(TestDBDsn)
	assertCanReadAndWriteToDisk(ctx, t, repo)

	assert.NoError(t, cont.Terminate(ctx))
}

func TestCanMakeQuery(t *testing.T) {
	ctx := context.Background()
	cont, err := createPsqlContainer(ctx)
	assert.NoError(t, err)

	repo := New(TestDBDsn)
	repo.Clear()
	repo.SetGauge(ctx, "g1", 1.1)
	repo.SetGauge(ctx, "g2", 2.2)

	repo.AddCounter(ctx, "c1", 1)
	repo.AddCounter(ctx, "c2", 2)

	cnt := repo.GetCounterOrZero(ctx, "c1")
	assert.Equal(t, int64(1), cnt)

	gg := repo.GetGaugeOrZero(ctx, "g1")
	assert.Equal(t, 1.1, gg)

	c0 := repo.GetCounterOrZero(ctx, "c0")
	assert.Equal(t, int64(0), c0)

	g0 := repo.GetGaugeOrZero(ctx, "g0")
	assert.Equal(t, float64(0), g0)

	cs := repo.GetAllCounters(ctx)
	assert.Equal(t, map[string]int64{"c1": int64(1), "c2": int64(2)}, cs)

	gs := repo.GetAllGauges(ctx)
	assert.Equal(t, map[string]float64{"g1": 1.1, "g2": 2.2}, gs)

	assert.Nil(t, repo.GetMemStorage())
	assert.NotNil(t, repo.GetDB())
	assert.NoError(t, repo.CloseDB())
	assert.NoError(t, cont.Terminate(ctx))
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
