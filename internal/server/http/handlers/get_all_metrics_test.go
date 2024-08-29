package handlers

import (
	"net/http"
	"testing"

	"github.com/gennadyterekhov/metrics-storage/internal/common/tests"

	"github.com/stretchr/testify/suite"

	"github.com/gennadyterekhov/metrics-storage/internal/common/testhelper"
	"github.com/stretchr/testify/assert"
)

type getAllTestSuite struct {
	tests.BaseSuiteWithServer
}

func (suite *getAllTestSuite) SetupSuite() {
	tests.InitBaseSuiteWithServer(suite)
}

func TestGetAll(t *testing.T) {
	suite.Run(t, new(getAllTestSuite))
}

func (suite *getAllTestSuite) TestGetAllMetrics() {
	type args struct {
		res http.ResponseWriter
		req *http.Request
	}
	cases := []struct {
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
	for _, tt := range cases {
		suite.T().Run(tt.name, func(t *testing.T) {
			response, responseBody := testhelper.SendRequest(
				t,
				suite.TestHTTPServer.Server,
				http.MethodGet,
				"/",
			)
			err := response.Body.Close()
			assert.NoError(t, err)

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

func (suite *getAllTestSuite) TestGetAllMetricsGzip() {
	type args struct {
		res http.ResponseWriter
		req *http.Request
	}
	cases := []struct {
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
	for _, tt := range cases {
		suite.T().Run(tt.name, func(t *testing.T) {
			response, responseBody := testhelper.SendGzipNoBodyRequest(
				t,
				suite.TestHTTPServer.Server,
				http.MethodGet,
				"/",
			)
			err := response.Body.Close()
			assert.NoError(t, err)

			assert.Equal(t,
				http.StatusOK,
				response.StatusCode,
			)
			assert.Equal(t,
				"text/html",
				response.Header.Get("Content-Type"),
			)
			assert.Equal(t,
				expected,
				string(responseBody),
			)
		})
	}
}
