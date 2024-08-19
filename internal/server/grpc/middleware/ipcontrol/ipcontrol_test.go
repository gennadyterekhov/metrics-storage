package ipcontrol

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"google.golang.org/grpc/metadata"

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

type grpcIPControlSuite struct {
	suite.Suite
	Server     *handlers.Server
	Client     pb.MetricsClient
	Conn       *grpc.ClientConn
	Repository *repositories.Repository
	Middleware *Middleware
}

func TestGrpcIPControlSuite(t *testing.T) {
	conf := config.New()
	DBOrRAM := storage.New(conf.DBDsn)
	repo := repositories.New(DBOrRAM)
	servicesPack := services.New(repo, conf)

	mdl := New(conf.TrustedSubnet)
	s := grpc.NewServer(
		grpc.UnaryInterceptor(mdl.IPControl),
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

	suiteInstance := &grpcIPControlSuite{}
	suiteInstance.Middleware = mdl
	suiteInstance.Repository = repo
	suiteInstance.Server = grpcServer
	c, conn := createClient(t)
	suiteInstance.Client = c
	suiteInstance.Conn = conn

	suite.Run(t, suiteInstance)
}

func (suite *grpcIPControlSuite) TearDownSuite() {
	err := suite.Conn.Close()
	assert.NoError(suite.T(), err)
}

func createClient(t *testing.T) (pb.MetricsClient, *grpc.ClientConn) {
	conn, err := grpc.NewClient("localhost:3333", grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err)
	return pb.NewMetricsClient(conn), conn
}

func (suite *grpcIPControlSuite) Test403IfSubnetDoesNotContainIP() {
	tests := []struct {
		ip            string
		trustedSubnet string
		code          int
	}{
		{
			ip:            "1.1.1.1",
			trustedSubnet: "1.1.1.0/24",
			code:          200,
		},
		{
			ip:            "2.1.1.1",
			trustedSubnet: "1.1.1.0/24",
			code:          403,
		},
		{
			ip:            "",
			trustedSubnet: "1.1.1.0/24",
			code:          403,
		},
		{
			ip:            "1.1.1.1",
			trustedSubnet: "",
			code:          200,
		},
		{
			ip:            "",
			trustedSubnet: "",
			code:          200,
		},
	}

	for i, tt := range tests {

		var err error
		_, ts, err := net.ParseCIDR(tt.trustedSubnet)
		if tt.trustedSubnet != "" {
			assert.NoError(suite.T(), err)
		}

		suite.T().Run(fmt.Sprintf("case%v", i), func(t *testing.T) {
			suite.Middleware.SetTrustedSubnet(ts)

			ctx := context.Background()
			ctx = metadata.AppendToOutgoingContext(ctx, "X-Real-IP", tt.ip)
			suite.Server.GetMetricService.Repository.AddCounter(ctx, "nm", 1)
			_, err := suite.Client.GetMetric(ctx, &pb.GetMetricRequest{Type: pb.MetricType_COUNTER, Name: "nm"})
			if tt.code == 200 {
				assert.NoError(suite.T(), err)
			}
			if tt.code != 200 {
				assert.Error(suite.T(), err)
			}
		})
	}
}

func (suite *grpcIPControlSuite) TestNilSubnet() {
	suite.Middleware.SetTrustedSubnet(nil)
	ctx := context.Background()

	ctx = metadata.AppendToOutgoingContext(ctx, "X-Real-IP", "1.1.1.1")
	suite.Server.GetMetricService.Repository.AddCounter(ctx, "nm", 1)
	_, err := suite.Client.GetMetric(ctx, &pb.GetMetricRequest{Type: pb.MetricType_COUNTER, Name: "nm"})

	assert.NoError(suite.T(), err)
}
