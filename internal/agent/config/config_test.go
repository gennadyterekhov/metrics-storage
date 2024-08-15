package config

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConfigValues(t *testing.T) {
	conf := New()

	assert.Equal(t, "localhost:8080", conf.Addr)
	assert.Equal(t, true, conf.IsGzip)
	assert.Equal(t, 10, conf.ReportInterval)
	assert.Equal(t, 2, conf.PollInterval)
	assert.Equal(t, false, conf.IsBatch)
	assert.Equal(t, "", conf.PayloadSignatureKey)
	assert.Equal(t, 5, conf.SimultaneousRequestsLimit)
	assert.Equal(t, "", conf.PublicKeyFilePath)
}

func TestCanGetConfigFromFile(t *testing.T) {
	exe, err := os.Getwd()
	assert.NoError(t, err)

	confPath := path.Join(path.Dir(exe), "config/testdata/config.json")
	err = os.Setenv("CONFIG", confPath)
	assert.NoError(t, err)

	conf := New()
	assert.Equal(t, 1, conf.ReportInterval)
	assert.Equal(t, "hello test", conf.Addr)
}

func TestEnvVarsOverwriteCliFlags(t *testing.T) {
	cmd := os.Args[0]
	os.Args = []string{cmd, "-r=2"}

	err := os.Setenv("REPORT_INTERVAL", "1")
	assert.NoError(t, err)

	conf := New()
	assert.Equal(t, 1, conf.ReportInterval)
}
