package handlers

import (
	"net/http"
	"testing"

	"github.com/gennadyterekhov/metrics-storage/internal/common/tests"
	"github.com/stretchr/testify/suite"

	"github.com/gennadyterekhov/metrics-storage/internal/common/testhelper"
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

func (suite *pingTestSuite) TestPing() {
	suite.T().Skip("only manual use because depends on host")
	type want struct {
		code int
	}
	type args struct {
		username string
	}

	cases := []struct {
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

	for _, tt := range cases {
		suite.T().Run(tt.name, func(t *testing.T) {
			assert.NoError(t, err)

			response, _ := testhelper.SendRequest(
				t,
				suite.TestHTTPServer.Server,
				http.MethodGet,
				"/ping",
			)
			err := response.Body.Close()
			if err != nil {
				panic(err)
			}
			assert.Equal(t, tt.want.code, response.StatusCode)
		})
	}
}
