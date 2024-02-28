package handlers

import (
	"github.com/gennadyterekhov/metrics-storage/internal/testhelper"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGetAllMetrics(t *testing.T) {
	type args struct {
		res http.ResponseWriter
		req *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "",
			args: args{},
		},
	}
	expected := `
<!DOCTYPE html>
<html>
  <head></head>
  <body>
    <h2>gauge</h2>
    <ul>

    </ul>
    <h2>counter</h2>
    <ul>

    </ul>
  </body>
</html>
`
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, responseBody := testhelper.SendRequest(
				t,
				testhelper.TestServer,
				http.MethodGet,
				"/",
			)

			assert.Equal(t,
				http.StatusOK,
				response.StatusCode,
			)
			assert.Equal(t,
				expected,
				string(responseBody),
			)

		})
	}
}
