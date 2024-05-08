package handlers

import (
	"testing"

	"github.com/gennadyterekhov/metrics-storage/internal/testhelper"
)

func TestMain(m *testing.M) {
	testhelper.BootstrapWithDefaultServer(m, GetRouter())
}
