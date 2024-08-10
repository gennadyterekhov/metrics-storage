package tests

import (
	"net/http/httptest"

	"github.com/stretchr/testify/suite"

	"github.com/gennadyterekhov/metrics-storage/internal/server/config"
	"github.com/gennadyterekhov/metrics-storage/internal/server/http/handlers/handlers"
	"github.com/gennadyterekhov/metrics-storage/internal/server/http/middleware"
	"github.com/gennadyterekhov/metrics-storage/internal/server/http/router"
	"github.com/gennadyterekhov/metrics-storage/internal/server/repositories"
	"github.com/gennadyterekhov/metrics-storage/internal/server/services/services"
	"github.com/gennadyterekhov/metrics-storage/internal/server/storage"
)

type TestHTTPServer struct {
	Server *httptest.Server
}
type HasRepo interface {
	SetRepository(repo *repositories.Repository)
	GetRepository() *repositories.Repository
}

type HasLifecycleMethods interface {
	SetupTest()
	TearDownTest()
}

type BaseSuiteInterface interface {
	HasLifecycleMethods
	HasRepo
}

type BaseSuite struct {
	suite.Suite
	Repository *repositories.Repository
}

func (suite *BaseSuite) SetupTest() {
	if suite.Repository != nil {
		suite.Repository.Clear()
	}
}

func (suite *BaseSuite) TearDownTest() {
	if suite.Repository != nil {
		suite.Repository.Clear()
	}
}

func InitBaseSuite[T BaseSuiteInterface](realSuite T) {
	repo := repositories.New(storage.New(""))
	realSuite.SetRepository(repo)
}

func (suite *BaseSuite) SetRepository(repo *repositories.Repository) {
	suite.Repository = repo
}

func (suite *BaseSuite) GetRepository() *repositories.Repository {
	return suite.Repository
}

type HasServer interface {
	SetServer(srv *httptest.Server)
	GetServer() *httptest.Server
}

type BaseSuiteWithServerInterface interface {
	BaseSuiteInterface
	HasServer
}

type BaseSuiteWithServer struct {
	BaseSuite
	TestHTTPServer
}

func InitBaseSuiteWithServer[T BaseSuiteWithServerInterface](srv T) {
	serverConfig := config.New()

	repo := repositories.New(storage.New(""))
	srv.SetRepository(repo)
	servs := services.New(repo, serverConfig)
	middlewareSet := middleware.New(serverConfig)
	controllersStruct := handlers.NewControllers(servs, middlewareSet)
	srv.SetServer(httptest.NewServer(
		router.New(controllersStruct).ChiRouter,
	))
}

func (s *BaseSuiteWithServer) SetRepository(repo *repositories.Repository) {
	s.Repository = repo
}

func (s *BaseSuiteWithServer) GetRepository() *repositories.Repository {
	return s.Repository
}

func (s *BaseSuiteWithServer) SetServer(srv *httptest.Server) {
	s.TestHTTPServer.Server = srv
}

func (s *BaseSuiteWithServer) GetServer() *httptest.Server {
	return s.TestHTTPServer.Server
}
