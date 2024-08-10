package config

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConfigValues(t *testing.T) {
	conf := New()

	assert.Equal(t, true, conf.Restore)
	assert.Equal(t, "<nil>", conf.TrustedSubnet.String())
	assert.Equal(t, "", conf.PrivateKeyFilePath)
	assert.Equal(t, "", conf.PayloadSignatureKey)
	assert.Equal(t, "", conf.DBDsn)
	assert.Equal(t, "/tmp/metrics-db.json", conf.FileStorage)
	assert.Equal(t, 300, conf.StoreInterval)
	assert.Equal(t, "localhost:8080", conf.Addr)
	assert.Equal(t, false, conf.IsGrpc)
}

func TestCanGetConfigFromFile(t *testing.T) {
	exe, err := os.Getwd()
	assert.NoError(t, err)

	confPath := path.Join(path.Dir(exe), "config/testdata/valid/config.json")
	err = os.Setenv("CONFIG", confPath)
	assert.NoError(t, err)

	conf := New()
	assert.Equal(t, 1, conf.StoreInterval)
	assert.Equal(t, "hello test", conf.Addr)
}

func TestCliFlagsOverwriteConfigFile(t *testing.T) {
	t.Skipf("cannot change flags in runtime")
	exe, err := os.Getwd()
	assert.NoError(t, err)

	confPath := path.Join(path.Dir(exe), "config/testdata/valid/config.json")
	err = os.Setenv("CONFIG", confPath)
	assert.NoError(t, err)

	cmd := os.Args[0]
	os.Args = []string{cmd, "-i=2"}

	conf := New()
	assert.Equal(t, 2, conf.StoreInterval)
}

func TestEnvVarsOverwriteCliFlags(t *testing.T) {
	cmd := os.Args[0]
	os.Args = []string{cmd, "-i=2"}

	err := os.Setenv("STORE_INTERVAL", "1")
	assert.NoError(t, err)

	conf := New()
	assert.Equal(t, 1, conf.StoreInterval)
}
