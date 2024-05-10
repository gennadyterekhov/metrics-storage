package handlers

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gennadyterekhov/metrics-storage/internal/server/config"
	"github.com/gennadyterekhov/metrics-storage/internal/server/storage"
	"github.com/gennadyterekhov/metrics-storage/internal/testhelper"
	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
	t.Skip("only manual use because depends on host")
	type want struct {
		code int
	}
	type args struct {
		username string
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "ok",
			args: args{username: "metrics_user"},
			want: want{code: http.StatusOK},
		},
		{
			name: "name unknown",
			args: args{username: "metrics_user2"},
			want: want{code: http.StatusInternalServerError},
		},
	}
	var err error

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			DBDsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
				`localhost`, tt.args.username, `metrics_pass`, `metrics_db_test`)

			config.Conf.DBDsn = DBDsn
			storage.MetricsRepository = storage.CreateDBStorage()
			assert.NoError(t, err)

			response, _ := testhelper.SendRequest(
				t,
				testhelper.TestServer,
				http.MethodGet,
				"/ping",
			)
			response.Body.Close()
			assert.Equal(t, tt.want.code, response.StatusCode)
		})
	}
}
