package tests

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gennadyterekhov/metrics-storage/internal/server/httpui/middleware"

	"github.com/gennadyterekhov/metrics-storage/internal/server/services/services"

	"github.com/gennadyterekhov/metrics-storage/internal/server/httpui/handlers/handlers"

	"github.com/gennadyterekhov/metrics-storage/internal/server/config"
	"github.com/gennadyterekhov/metrics-storage/internal/server/httpui/router"
	"github.com/gennadyterekhov/metrics-storage/internal/server/repositories"
	"github.com/gennadyterekhov/metrics-storage/internal/server/storage"
	"github.com/stretchr/testify/suite"
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
	realSuite.SetRepository(&repo)
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
	srv.SetRepository(&repo)
	servs := services.New(repo, serverConfig)
	middlewareSet := middleware.New(serverConfig)
	controllersStruct := handlers.NewControllers(&servs, middlewareSet)
	srv.SetServer(httptest.NewServer(
		router.New(&controllersStruct).ChiRouter,
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

func (ts *TestHTTPServer) SendGet(
	path string,
	token string,
) (int, []byte) {
	req, err := http.NewRequest(http.MethodGet, ts.Server.URL+path, strings.NewReader(""))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	response, err := ts.Server.Client().Do(req)
	if err != nil {
		panic(err)
	}
	bodyAsBytes, err := getBodyAsBytes(response.Body)
	if err != nil {
		return 0, nil
	}
	err = response.Body.Close()
	if err != nil {
		panic(err)
	}
	return response.StatusCode, bodyAsBytes
}

func (ts *TestHTTPServer) SendPostWithoutToken(
	path string,
	requestBody *bytes.Buffer,
) int {
	code, _ := ts.SendPostAndReturnBody(path, "application/json", "", requestBody)

	return code
}

func (ts *TestHTTPServer) SendPost(
	path string,
	contentType string,
	token string,
	requestBody *bytes.Buffer,
) int {
	code, _ := ts.SendPostAndReturnBody(path, contentType, token, requestBody)

	return code
}

func (ts *TestHTTPServer) SendPostAndReturnBody(
	path string,
	contentType string,
	token string,
	requestBody *bytes.Buffer,
) (int, []byte) {
	req, err := http.NewRequest(http.MethodPost, ts.Server.URL+path, requestBody)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", token)

	response, err := ts.Server.Client().Do(req)
	if err != nil {
		panic(err)
	}
	bodyAsBytes, err := getBodyAsBytes(response.Body)
	if err != nil {
		panic(err)
	}
	err = response.Body.Close()
	if err != nil {
		panic(err)
	}
	return response.StatusCode, bodyAsBytes
}

func getBodyAsBytes(reader io.Reader) ([]byte, error) {
	readBytes, err := io.ReadAll(reader)
	if err != nil {
		return []byte{}, err
	}

	return readBytes, nil
}
