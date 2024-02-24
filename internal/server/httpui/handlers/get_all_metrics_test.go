package handlers

import (
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			GetAllMetrics(w, tt.args.req)

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

			res := w.Result()
			metricFromResponse, _ := io.ReadAll(res.Body)
			defer res.Body.Close()
			assert.Equal(t,
				expected,
				string(metricFromResponse),
			)

		})
	}
}
