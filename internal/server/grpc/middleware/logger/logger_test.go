package logger

import (
	"context"
	"net"
	"testing"
	"time"

	pb "github.com/gennadyterekhov/metrics-storage/internal/common/protobuf"
	"github.com/gennadyterekhov/metrics-storage/internal/server/config"
	"github.com/gennadyterekhov/metrics-storage/internal/server/grpc/handlers"
	"github.com/gennadyterekhov/metrics-storage/internal/server/repositories"
	"github.com/gennadyterekhov/metrics-storage/internal/server/services/services"
	"github.com/gennadyterekhov/metrics-storage/internal/server/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type grpcMiddlewareSuite struct {
	suite.Suite
	Server     *handlers.Server
	Client     pb.MetricsClient
	Conn       *grpc.ClientConn
	Repository *repositories.Repository
}

func TestGrpcMiddlewareSuite(t *testing.T) {
	conf := config.New()
	DBOrRAM := storage.New(conf.DBDsn)
	repo := repositories.New(DBOrRAM)
	servicesPack := services.New(repo, conf)

	s := grpc.NewServer(
		grpc.UnaryInterceptor(LoggingInterceptor),
	)
	grpcServer := &handlers.Server{
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

	suiteInstance := &grpcMiddlewareSuite{}
	suiteInstance.Repository = repo
	suiteInstance.Server = grpcServer
	c, conn := createClient(t)
	suiteInstance.Client = c
	suiteInstance.Conn = conn

	suite.Run(t, suiteInstance)
}

func (suite *grpcMiddlewareSuite) TearDownSuite() {
	err := suite.Conn.Close()
	assert.NoError(suite.T(), err)
}

func (suite *grpcMiddlewareSuite) TestGrpcMiddleware() {
	var err error
	ctx := context.Background()

	suite.Server.GetMetricService.Repository.AddCounter(ctx, "nm", 1)
	_, err = suite.Client.GetMetric(ctx, &pb.GetMetricRequest{Type: pb.MetricType_COUNTER, Name: "nm"})
	assert.NoError(suite.T(), err)
}

func createClient(t *testing.T) (pb.MetricsClient, *grpc.ClientConn) {
	conn, err := grpc.NewClient("localhost:3333", grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err)
	return pb.NewMetricsClient(conn), conn
}
