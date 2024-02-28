package handlers

import (
	"github.com/gennadyterekhov/metrics-storage/internal/testhelper"
	"net/http/httptest"
	"testing"
)

func TestMain(m *testing.M) {
	testhelper.BootstrapWithServer(
		m,
		httptest.NewServer(
			GetRouter(),
		),
	)
}
