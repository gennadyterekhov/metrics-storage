// because the source for agent and server are in the same internal,
// we can use server's code without actually launching the server's binary and making requests
// so, this file uses httptest.Server
package agent

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gennadyterekhov/metrics-storage/internal/agent/client"
	"github.com/gennadyterekhov/metrics-storage/internal/agent/poller"
	"github.com/gennadyterekhov/metrics-storage/internal/agent/sender"

	agentConfig "github.com/gennadyterekhov/metrics-storage/internal/agent/config"
	"github.com/gennadyterekhov/metrics-storage/internal/common/testhelper"
	"github.com/gennadyterekhov/metrics-storage/internal/common/tests"
	"github.com/gennadyterekhov/metrics-storage/internal/server/config"
	"github.com/gennadyterekhov/metrics-storage/internal/server/httpui/handlers/handlers"
	"github.com/gennadyterekhov/metrics-storage/internal/server/httpui/middleware"
	"github.com/gennadyterekhov/metrics-storage/internal/server/httpui/router"
	"github.com/gennadyterekhov/metrics-storage/internal/server/repositories"
	"github.com/gennadyterekhov/metrics-storage/internal/server/services/services"
	"github.com/gennadyterekhov/metrics-storage/internal/server/storage"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type agentTestSuite struct {
	tests.BaseSuiteWithServer
}

func NewWithCustomConfig(conf *agentConfig.Config) *Agent {
	metricsStorageClient := client.New(conf)

	inst := &Agent{
		Config:               conf,
		Poller:               poller.New(conf.PollInterval),
		MetricsStorageClient: metricsStorageClient,
		Sender:               sender.New(metricsStorageClient, conf),
	}

	return inst
}

func (suite *agentTestSuite) SetupSuite() {
	tests.InitBaseSuiteWithServer(suite)
}

func TestAgentSuite(t *testing.T) {
	suite.Run(t, new(agentTestSuite))
}

func (suite *agentTestSuite) TestAgent() {
	ctx, cancelContextFn := context.WithTimeout(context.Background(), 300*time.Millisecond)

	defer cancelContextFn()

	go runAgentRoutine(ctx, &agentConfig.Config{
		Addr:                      suite.TestHTTPServer.Server.URL,
		ReportInterval:            1,
		PollInterval:              1,
		SimultaneousRequestsLimit: 5,
	})

	<-ctx.Done()

	contextEndCondition := ctx.Err()

	if contextEndCondition == context.DeadlineExceeded || contextEndCondition == context.Canceled {
		totalCounters := len(suite.Repository.GetAllCounters(context.Background()))
		totalGauges := len(suite.Repository.GetAllGauges(context.Background()))

		assert.Equal(suite.T(),
			1,
			totalCounters,
		)
		assert.LessOrEqual(suite.T(),
			27+1,
			totalGauges,
		)
	} else {
		suite.T().Error("context didnt finish")
	}
}

func (suite *agentTestSuite) TestList() {
	ctx, cancelContextFn := context.WithTimeout(context.Background(), 300*time.Millisecond)

	defer cancelContextFn()

	go runAgentRoutine(ctx, &agentConfig.Config{
		Addr:                      suite.TestHTTPServer.Server.URL,
		ReportInterval:            1,
		PollInterval:              1,
		IsBatch:                   true,
		SimultaneousRequestsLimit: 5,
	})

	<-ctx.Done()

	contextEndCondition := ctx.Err()

	if contextEndCondition == context.DeadlineExceeded || contextEndCondition == context.Canceled {
		assert.Equal(suite.T(),
			1,
			len(suite.Repository.GetAllCounters(context.Background())),
		)
		assert.LessOrEqual(suite.T(),
			27+1,
			len(suite.Repository.GetAllGauges(context.Background())),
		)

		return
	}
	suite.T().Error("context didnt finish")
}

func (suite *agentTestSuite) TestGzip() {
	ctx, cancelContextFn := context.WithTimeout(context.Background(), 300*time.Millisecond)

	defer cancelContextFn()
	go runAgentRoutine(ctx, &agentConfig.Config{
		Addr:                      suite.TestHTTPServer.Server.URL,
		ReportInterval:            1,
		PollInterval:              1,
		IsGzip:                    true,
		SimultaneousRequestsLimit: 5,
	})

	<-ctx.Done()

	contextEndCondition := ctx.Err()

	if contextEndCondition == context.DeadlineExceeded || contextEndCondition == context.Canceled {
		assert.Equal(suite.T(),
			1,
			len(suite.Repository.GetAllCounters(context.Background())),
		)
		assert.LessOrEqual(suite.T(),
			27+1,
			len(suite.Repository.GetAllGauges(context.Background())),
		)
		savedValue := suite.Repository.GetCounterOrZero(context.Background(), "PollCount")
		assert.Equal(suite.T(), int64(1), savedValue)
		return
	}

	suite.T().Error("context didnt finish")
}

func (suite *agentTestSuite) TestSameValueReturnedFromServer() {
	ctx, cancelContextFn := context.WithTimeout(context.Background(), 300*time.Millisecond)

	defer cancelContextFn()
	go runAgentRoutine(ctx, &agentConfig.Config{
		Addr:                      suite.TestHTTPServer.Server.URL,
		ReportInterval:            1,
		PollInterval:              1,
		IsBatch:                   true,
		SimultaneousRequestsLimit: 5,
	})

	<-ctx.Done()
	contextEndCondition := ctx.Err()

	if contextEndCondition == context.DeadlineExceeded || contextEndCondition == context.Canceled {
		assert.Equal(suite.T(),
			1,
			len(suite.Repository.GetAllCounters(context.Background())),
		)
		assert.LessOrEqual(suite.T(),
			27+1,
			len(suite.Repository.GetAllGauges(context.Background())),
		)

		url := "/value/gauge/BuckHashSys"

		r, responseBody := testhelper.SendRequest(suite.T(), suite.TestHTTPServer.Server, http.MethodGet, url)
		err := r.Body.Close()
		if err != nil {
			return
		}
		savedValue := suite.Repository.GetGaugeOrZero(context.Background(), "BuckHashSys")

		assert.Equal(
			suite.T(),
			strconv.FormatFloat(savedValue, 'g', -1, 64),
			responseBody,
		)
	} else {
		suite.T().Error("context didnt finish")
	}
}

func (suite *agentTestSuite) TestReportIntervalMoreThanPollInterval() {
	ctx, cancelContextFn := context.WithTimeout(context.Background(), 3000*time.Millisecond)
	defer cancelContextFn()

	go runAgentRoutine(ctx, &agentConfig.Config{
		Addr:                      suite.TestHTTPServer.Server.URL,
		ReportInterval:            2,
		PollInterval:              1,
		IsBatch:                   true,
		SimultaneousRequestsLimit: 5,
	})

	<-ctx.Done()
	contextEndCondition := ctx.Err()

	if contextEndCondition == context.DeadlineExceeded || contextEndCondition == context.Canceled {
		assert.Equal(suite.T(),
			1,
			len(suite.Repository.GetAllCounters(context.Background())),
		)
		assert.LessOrEqual(suite.T(),
			27+1,
			len(suite.Repository.GetAllGauges(context.Background())),
		)
	} else {
		suite.T().Error("context didnt finish")
	}
}

func (suite *agentTestSuite) TestReportIntervalLessThanPollInterval() {
	ctx, cancelContextFn := context.WithTimeout(context.Background(), 3000*time.Millisecond)
	defer cancelContextFn()

	go runAgentRoutine(ctx, &agentConfig.Config{
		Addr:                      suite.TestHTTPServer.Server.URL,
		ReportInterval:            1,
		PollInterval:              2,
		IsBatch:                   true,
		SimultaneousRequestsLimit: 5,
	})

	<-ctx.Done()
	contextEndCondition := ctx.Err()

	if contextEndCondition == context.DeadlineExceeded || contextEndCondition == context.Canceled {
		assert.Equal(suite.T(),
			1,
			len(suite.Repository.GetAllCounters(context.Background())),
		)
		assert.LessOrEqual(suite.T(),
			27+1,
			len(suite.Repository.GetAllGauges(context.Background())),
		)
	} else {
		suite.T().Error("context didnt finish")
	}
}

// TestAsymmetricEncryptionUsingKeyFiles is not part of suite because it relies on overridden config
func TestAsymmetricEncryptionUsingKeyFiles(t *testing.T) {
	type app struct {
		tests.BaseSuiteWithServer
	}
	appCustomServer := app{}

	privateKeyFilePath := "../../keys/private.test"
	serverConfig := config.ServerConfig{
		Addr:                "",
		StoreInterval:       0,
		FileStorage:         "",
		Restore:             false,
		DBDsn:               "",
		PayloadSignatureKey: "",
		PrivateKeyFilePath:  privateKeyFilePath,
	}

	repo := repositories.New(storage.New(""))
	appCustomServer.SetRepository(&repo)
	servs := services.New(repo, &serverConfig)
	middlewareSet := middleware.New(&serverConfig)
	controllersStruct := handlers.NewControllers(&servs, middlewareSet)
	testServer := httptest.NewServer(
		router.New(&controllersStruct).ChiRouter,
	)
	appCustomServer.SetServer(testServer)

	ctx, cancelContextFn := context.WithTimeout(context.Background(), 300*time.Millisecond)
	publicKeyFilePath := "../../keys/public.test"
	defer cancelContextFn()
	go runAgentRoutine(ctx, &agentConfig.Config{
		Addr:                      appCustomServer.TestHTTPServer.Server.URL,
		ReportInterval:            1,
		PollInterval:              1,
		IsGzip:                    true,
		SimultaneousRequestsLimit: 5,
		PublicKeyFilePath:         publicKeyFilePath,
	})

	<-ctx.Done()

	contextEndCondition := ctx.Err()

	if contextEndCondition == context.DeadlineExceeded || contextEndCondition == context.Canceled {
		assert.Equal(t,
			1,
			len(appCustomServer.Repository.GetAllCounters(context.Background())),
		)
		assert.LessOrEqual(t,
			27+1,
			len(appCustomServer.Repository.GetAllGauges(context.Background())),
		)
		savedValue := appCustomServer.Repository.GetCounterOrZero(context.Background(), "PollCount")
		assert.Equal(t, int64(1), savedValue)
		return
	}
	t.Error("context didnt finish")
}

func runAgentRoutine(ctx context.Context, config *agentConfig.Config) {
	instance := NewWithCustomConfig(config)
	instance.RunAgent(ctx)
}
