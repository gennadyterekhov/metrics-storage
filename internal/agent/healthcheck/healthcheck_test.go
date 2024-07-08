package healthcheck

import (
	"net/http"
	"testing"

	"github.com/gennadyterekhov/metrics-storage/internal/common/tests"
	"github.com/stretchr/testify/suite"

	"github.com/stretchr/testify/assert"
)

type healthcheckTestSuite struct {
	tests.BaseSuiteWithServer
}

func (suite *healthcheckTestSuite) SetupSuite() {
	tests.InitBaseSuiteWithServer(suite)
}

func TestHealthcheck(t *testing.T) {
	suite.Run(t, new(healthcheckTestSuite))
}

func (st *healthcheckTestSuite) TestCanSendHeadRequest() {
	type want struct {
		code int
	}
	cases := []struct {
		name string
		url  string
		want want
	}{
		{
			name: "head",
			url:  "/",
			want: want{code: http.StatusOK},
		},
	}

	for _, tt := range cases {
		st.T().Run(tt.name, func(t *testing.T) {
			isOk := MakeHealthcheck(st.TestHTTPServer.Server.URL)
			assert.Equal(t, true, isOk)
		})
	}
}
