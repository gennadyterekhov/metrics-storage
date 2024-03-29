package handlers

import (
	"github.com/gennadyterekhov/metrics-storage/internal/testhelper"
	"testing"
)

func TestMain(m *testing.M) {
	testhelper.BootstrapWithDefaultServer(m, GetRouter())

}
