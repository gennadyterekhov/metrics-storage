package handlers

import (
	"testing"

	"github.com/gennadyterekhov/metrics-storage/internal/server/httpui/router"

	"github.com/gennadyterekhov/metrics-storage/internal/common/testhelper"
)

func TestMain(m *testing.M) {
	testhelper.BootstrapWithDefaultServer(m, router.GetRouter())
}
