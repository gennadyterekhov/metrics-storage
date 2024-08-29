package handlers

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	pb "github.com/gennadyterekhov/metrics-storage/internal/common/protobuf"
	"github.com/gennadyterekhov/metrics-storage/internal/server/config"
	"github.com/gennadyterekhov/metrics-storage/internal/server/repositories"
	"github.com/gennadyterekhov/metrics-storage/internal/server/services/services"
	"github.com/gennadyterekhov/metrics-storage/internal/server/storage"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type grpcHandlersSuite struct {
	suite.Suite
	Server     *Server
	Client     pb.MetricsClient
	Conn       *grpc.ClientConn
	Repository *repositories.Repository
}

func TestGrpcServerSuite(t *testing.T) {
	conf := config.New()
	DBOrRAM := storage.New(conf.DBDsn)
	repo := repositories.New(DBOrRAM)
	servicesPack := services.New(repo, conf)

	s := grpc.NewServer()
	grpcServer := &Server{
		GetMetricService:  servicesPack.GetMetricService,
		SaveMetricService: servicesPack.SaveMetricService,
		PingService:       servicesPack.PingService,
	}
	pb.RegisterMetricsServer(s, grpcServer)
	listen, err := net.Listen("tcp", "localhost:3333")
	assert.NoError(nil, err)

	go func() {
		err = s.Serve(listen)
		assert.NoError(nil, err)
	}()
	time.Sleep(10 * time.Millisecond)

	suiteInstance := &grpcHandlersSuite{}
	suiteInstance.Repository = repo
	suiteInstance.Server = grpcServer
	c, conn := createClient(t)
	suiteInstance.Client = c
	suiteInstance.Conn = conn

	suite.Run(t, suiteInstance)
}

func (suite *grpcHandlersSuite) SetupTest() {
	if suite.Repository != nil {
		suite.Repository.Clear()
	}
}

func (suite *grpcHandlersSuite) TearDownSuite() {
	err := suite.Conn.Close()
	assert.NoError(suite.T(), err)
}

func (suite *grpcHandlersSuite) TestPingGrpc() {
	_, err := suite.Client.Ping(context.Background(), &pb.PingRequest{})
	assert.NoError(suite.T(), err)
}

func (suite *grpcHandlersSuite) TestGetMetricGrpc() {
	var err error
	ctx := context.Background()

	suite.Server.GetMetricService.Repository.AddCounter(ctx, "nm", 1)
	resp, err := suite.Client.GetMetric(ctx, &pb.GetMetricRequest{Type: pb.MetricType_COUNTER, Name: "nm"})
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(1), resp.Counter)
	assert.Equal(suite.T(), pb.MetricType_COUNTER, resp.Type)

	suite.Server.GetMetricService.Repository.AddCounter(ctx, "nm", 2)
	resp, err = suite.Client.GetMetric(ctx, &pb.GetMetricRequest{Type: pb.MetricType_COUNTER, Name: "nm"})
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(1+2), resp.Counter)
	assert.Equal(suite.T(), pb.MetricType_COUNTER, resp.Type)

	suite.Server.GetMetricService.Repository.SetGauge(ctx, "gg1", 1.1)
	resp, err = suite.Client.GetMetric(ctx, &pb.GetMetricRequest{Type: pb.MetricType_GAUGE, Name: "gg1"})
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1.1, resp.Gauge)
	assert.Equal(suite.T(), pb.MetricType_GAUGE, resp.Type)

	suite.Server.GetMetricService.Repository.SetGauge(ctx, "gg1", 2.2)
	resp, err = suite.Client.GetMetric(ctx, &pb.GetMetricRequest{Type: pb.MetricType_GAUGE, Name: "gg1"})
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 2.2, resp.Gauge)
	assert.Equal(suite.T(), pb.MetricType_GAUGE, resp.Type)
}

func (suite *grpcHandlersSuite) TestGetAllMetricGrpc() {
	var err error
	ctx := context.Background()

	suite.Server.GetMetricService.Repository.AddCounter(ctx, "c1", 1)
	suite.Server.GetMetricService.Repository.AddCounter(ctx, "c2", 2)
	suite.Server.GetMetricService.Repository.SetGauge(ctx, "g1", 1.1)
	suite.Server.GetMetricService.Repository.SetGauge(ctx, "g2", 2.2)

	resp, err := suite.Client.GetAllMetrics(ctx, &pb.GetAllMetricsRequest{})
	assert.NoError(suite.T(), err)
	assert.NotEqual(suite.T(), "", resp.Html) // cannot test real string because order can be different
}

func (suite *grpcHandlersSuite) TestSaveMetricGrpc() {
	var err error
	ctx := context.Background()

	_, err = suite.Client.SaveMetric(ctx, &pb.SaveMetricRequest{Type: pb.MetricType_COUNTER, Name: "c", Counter: 1})
	assert.NoError(suite.T(), err)

	_, err = suite.Client.SaveMetric(ctx, &pb.SaveMetricRequest{Type: pb.MetricType_COUNTER, Name: "c", Counter: 2})
	assert.NoError(suite.T(), err)

	respOfGet, err := suite.Client.GetMetric(ctx, &pb.GetMetricRequest{Type: pb.MetricType_COUNTER, Name: "c"})
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(1+2), respOfGet.Counter)

	_, err = suite.Client.SaveMetric(ctx, &pb.SaveMetricRequest{Type: pb.MetricType_GAUGE, Name: "g", Gauge: 1.1})
	assert.NoError(suite.T(), err)

	_, err = suite.Client.SaveMetric(ctx, &pb.SaveMetricRequest{Type: pb.MetricType_GAUGE, Name: "g", Gauge: 2.2})
	assert.NoError(suite.T(), err)

	respOfGet, err = suite.Client.GetMetric(ctx, &pb.GetMetricRequest{Type: pb.MetricType_GAUGE, Name: "g"})
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 2.2, respOfGet.Gauge)
}

func (suite *grpcHandlersSuite) TestSaveMetricListGrpc() {
	var err error
	ctx := context.Background()

	c1 := &pb.SaveMetricRequest{Type: pb.MetricType_COUNTER, Name: "c", Counter: 1}
	g1 := &pb.SaveMetricRequest{Type: pb.MetricType_GAUGE, Name: "g1", Gauge: 1.1}
	g2 := &pb.SaveMetricRequest{Type: pb.MetricType_GAUGE, Name: "g2", Gauge: 2.2}
	req := &pb.SaveMetricListRequest{Request: []*pb.SaveMetricRequest{c1, g1, g2}}
	_, err = suite.Client.SaveMetricList(ctx, req)
	assert.NoError(suite.T(), err)

	cnt, err := suite.Server.SaveMetricService.Repository.GetCounter(ctx, "c")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(1), cnt)

	g1v, err := suite.Server.SaveMetricService.Repository.GetGauge(ctx, "g1")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, 1, g1v)

	g2v, err := suite.Server.SaveMetricService.Repository.GetGauge(ctx, "g2")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 2.2, g2v)
}

func createClient(t *testing.T) (pb.MetricsClient, *grpc.ClientConn) {
	conn, err := grpc.NewClient("localhost:3333", grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err)
	return pb.NewMetricsClient(conn), conn
}
