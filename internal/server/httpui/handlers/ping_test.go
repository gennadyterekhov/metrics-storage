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
	var err error

	assert.NoError(suite.T(), err)

	response, _ := testhelper.SendRequest(
		suite.T(),
		suite.TestHTTPServer.Server,
		http.MethodGet,
		"/ping",
	)
	err = response.Body.Close()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, response.StatusCode)
}
