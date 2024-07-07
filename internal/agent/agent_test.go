// because the source for agent and server are in the same internal,
// we can use server's code without actually launching the server's binary and making requests
// so, this file uses httptest.Server
package agent

import (
	"context"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/gennadyterekhov/metrics-storage/internal/common/tests"
	"github.com/stretchr/testify/suite"

	"github.com/gennadyterekhov/metrics-storage/internal/common/testhelper"
	"github.com/stretchr/testify/assert"
)

type agentTestSuite struct {
	tests.BaseSuiteWithServer
}

func (suite *agentTestSuite) SetupSuite() {
	tests.InitBaseSuiteWithServer(suite)
}

func TestAgentSuite(t *testing.T) {
	suite.Run(t, new(agentTestSuite))
}

func (st *agentTestSuite) TestAgent() {
	ctx, cancelContextFn := context.WithTimeout(context.Background(), 300*time.Millisecond)

	defer cancelContextFn()

	go runAgentRoutine(ctx, &AgentConfig{
		Addr:                      st.TestHTTPServer.Server.URL, //
		ReportInterval:            1,
		PollInterval:              1,
		SimultaneousRequestsLimit: 5,
	})

	<-ctx.Done()

	contextEndCondition := ctx.Err()

	if contextEndCondition == context.DeadlineExceeded || contextEndCondition == context.Canceled {
		totalCounters := len(st.Repository.GetAllCounters(context.Background()))
		totalGauges := len(st.Repository.GetAllGauges(context.Background()))

		assert.Equal(st.T(),
			1,
			totalCounters,
		)
		assert.LessOrEqual(st.T(),
			27+1,
			totalGauges,
		)
	} else {
		st.T().Error("context didnt finish")
	}
}

func (st *agentTestSuite) TestList() {
	ctx, cancelContextFn := context.WithTimeout(context.Background(), 300*time.Millisecond)

	defer cancelContextFn()

	go runAgentRoutine(ctx, &AgentConfig{
		Addr:                      st.TestHTTPServer.Server.URL,
		ReportInterval:            1,
		PollInterval:              1,
		IsBatch:                   true,
		SimultaneousRequestsLimit: 5,
	})

	<-ctx.Done()

	contextEndCondition := ctx.Err()

	if contextEndCondition == context.DeadlineExceeded || contextEndCondition == context.Canceled {
		assert.Equal(st.T(),
			1,
			len(st.Repository.GetAllCounters(context.Background())),
		)
		assert.LessOrEqual(st.T(),
			27+1,
			len(st.Repository.GetAllGauges(context.Background())),
		)

		return
	} else {
		st.T().Error("context didnt finish")
	}
}

func (st *agentTestSuite) TestGzip() {
	ctx, cancelContextFn := context.WithTimeout(context.Background(), 300*time.Millisecond)

	defer cancelContextFn()
	go runAgentRoutine(ctx, &AgentConfig{
		Addr:                      st.TestHTTPServer.Server.URL,
		ReportInterval:            1,
		PollInterval:              1,
		IsGzip:                    true,
		SimultaneousRequestsLimit: 5,
	})

	<-ctx.Done()

	contextEndCondition := ctx.Err()

	if contextEndCondition == context.DeadlineExceeded || contextEndCondition == context.Canceled {
		assert.Equal(st.T(),
			1,
			len(st.Repository.GetAllCounters(context.Background())),
		)
		assert.LessOrEqual(st.T(),
			27+1,
			len(st.Repository.GetAllGauges(context.Background())),
		)
		savedValue := st.Repository.GetCounterOrZero(context.Background(), "PollCount")
		assert.Equal(st.T(), int64(1), savedValue)
		return
	}

	st.T().Error("context didnt finish")
}

func (st *agentTestSuite) TestSameValueReturnedFromServer() {
	ctx, cancelContextFn := context.WithTimeout(context.Background(), 300*time.Millisecond)

	defer cancelContextFn()
	go runAgentRoutine(ctx, &AgentConfig{
		Addr:                      st.TestHTTPServer.Server.URL,
		ReportInterval:            1,
		PollInterval:              1,
		IsBatch:                   true,
		SimultaneousRequestsLimit: 5,
	})

	<-ctx.Done()
	contextEndCondition := ctx.Err()

	if contextEndCondition == context.DeadlineExceeded || contextEndCondition == context.Canceled {
		assert.Equal(st.T(),
			1,
			len(st.Repository.GetAllCounters(context.Background())),
		)
		assert.LessOrEqual(st.T(),
			27+1,
			len(st.Repository.GetAllGauges(context.Background())),
		)

		url := "/value/gauge/BuckHashSys"

		r, responseBody := testhelper.SendRequest(st.T(), st.TestHTTPServer.Server, http.MethodGet, url)
		r.Body.Close()
		savedValue := st.Repository.GetGaugeOrZero(context.Background(), "BuckHashSys")

		assert.Equal(
			st.T(),
			strconv.FormatFloat(savedValue, 'g', -1, 64),
			string(responseBody),
		)
	} else {
		st.T().Error("context didnt finish")
	}
}

func (st *agentTestSuite) TestReportIntervalMoreThanPollInterval() {
	ctx, cancelContextFn := context.WithTimeout(context.Background(), 3000*time.Millisecond)
	defer cancelContextFn()

	go runAgentRoutine(ctx, &AgentConfig{
		Addr:                      st.TestHTTPServer.Server.URL,
		ReportInterval:            2,
		PollInterval:              1,
		IsBatch:                   true,
		SimultaneousRequestsLimit: 5,
	})

	<-ctx.Done()
	contextEndCondition := ctx.Err()

	if contextEndCondition == context.DeadlineExceeded || contextEndCondition == context.Canceled {
		assert.Equal(st.T(),
			1,
			len(st.Repository.GetAllCounters(context.Background())),
		)
		assert.LessOrEqual(st.T(),
			27+1,
			len(st.Repository.GetAllGauges(context.Background())),
		)
	} else {
		st.T().Error("context didnt finish")
	}
}

func (st *agentTestSuite) TestReportIntervalLessThanPollInterval() {
	ctx, cancelContextFn := context.WithTimeout(context.Background(), 3000*time.Millisecond)
	defer cancelContextFn()

	go runAgentRoutine(ctx, &AgentConfig{
		Addr:                      st.TestHTTPServer.Server.URL,
		ReportInterval:            1,
		PollInterval:              2,
		IsBatch:                   true,
		SimultaneousRequestsLimit: 5,
	})

	<-ctx.Done()
	contextEndCondition := ctx.Err()

	if contextEndCondition == context.DeadlineExceeded || contextEndCondition == context.Canceled {
		assert.Equal(st.T(),
			1,
			len(st.Repository.GetAllCounters(context.Background())),
		)
		assert.LessOrEqual(st.T(),
			27+1,
			len(st.Repository.GetAllGauges(context.Background())),
		)
	} else {
		st.T().Error("context didnt finish")
	}
}

func runAgentRoutine(ctx context.Context, config *AgentConfig) {
	RunAgent(ctx, config)
}
