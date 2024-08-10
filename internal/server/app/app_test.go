package app

import (
	"context"
	"os"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestApp(t *testing.T) {
	prepareStoragePath(t)

	app := New()

	assert.NotEqual(t, nil, app.Config)
	assert.NotEqual(t, nil, app.Repository)
	assert.NotEqual(t, nil, app.Services)
	assert.NotEqual(t, nil, app.Controllers)
	assert.NotEqual(t, nil, app.DBOrRAM)
	assert.NotEqual(t, nil, app.Router)
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	go func() {
		err := app.StartServer(ctx)
		assert.Error(t, err)
		assert.Equal(t, "http: Server closed", err.Error())
	}()
	time.Sleep(70 * time.Millisecond)

	cancel()
}

func TestAppGrpc(t *testing.T) {
	prepareStoragePath(t)
	prepareGrpc(t)

	app := New()

	assert.NotEqual(t, nil, app.Config)
	assert.NotEqual(t, nil, app.Repository)
	assert.NotEqual(t, nil, app.Services)
	assert.NotEqual(t, nil, app.Controllers)
	assert.NotEqual(t, nil, app.DBOrRAM)
	assert.NotEqual(t, nil, app.Router)
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	go func() {
		err := app.StartServer(ctx)
		assert.NoError(t, err)
	}()
	time.Sleep(70 * time.Millisecond)

	cancel()
}

func prepareStoragePath(t *testing.T) {
	exe, err := os.Getwd()
	assert.NoError(t, err)

	storagePath := path.Join(path.Dir(exe), "app/testdata/metrics.json")
	err = os.Setenv("FILE_STORAGE_PATH", storagePath)

	assert.NoError(t, err)
}

func prepareGrpc(t *testing.T) {
	err := os.Setenv("USE_GRPC", "1")
	assert.NoError(t, err)
}
