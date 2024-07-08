package handlers

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gennadyterekhov/metrics-storage/internal/common/tests"
	"github.com/stretchr/testify/suite"

	"github.com/gennadyterekhov/metrics-storage/internal/common/testhelper"
	"github.com/gennadyterekhov/metrics-storage/internal/server/config"
	"github.com/stretchr/testify/assert"
)

type pingTestSuite struct {
	tests.BaseSuiteWithServer
}

func (suite *pingTestSuite) SetupSuite() {
	tests.InitBaseSuiteWithServer(suite)
}

func TestPingHandler(t *testing.T) {
	suite.Run(t, new(pingTestSuite))
}

func (st *pingTestSuite) TestPing() {
	st.T().Skip("only manual use because depends on host")
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
		st.T().Run(tt.name, func(t *testing.T) {
			DBDsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
				`localhost`, tt.args.username, `metrics_pass`, `metrics_db_test`)

			config.Conf.DBDsn = DBDsn
			assert.NoError(t, err)

			response, _ := testhelper.SendRequest(
				t,
				st.TestHTTPServer.Server,
				http.MethodGet,
				"/ping",
			)
			response.Body.Close()
			assert.Equal(t, tt.want.code, response.StatusCode)
		})
	}
}
