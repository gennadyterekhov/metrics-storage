package app

import (
	"context"
	"net"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/gennadyterekhov/metrics-storage/internal/server/grpc/middleware/ipcontrol"

	"google.golang.org/grpc"

	"github.com/gennadyterekhov/metrics-storage/internal/common/logger"
	pb "github.com/gennadyterekhov/metrics-storage/internal/common/protobuf"
	"github.com/gennadyterekhov/metrics-storage/internal/server/config"
	grpcHandlers "github.com/gennadyterekhov/metrics-storage/internal/server/grpc/handlers"
	loggerInterceptor "github.com/gennadyterekhov/metrics-storage/internal/server/grpc/middleware/logger"
	"github.com/gennadyterekhov/metrics-storage/internal/server/http/handlers/handlers"
	"github.com/gennadyterekhov/metrics-storage/internal/server/http/middleware"
	"github.com/gennadyterekhov/metrics-storage/internal/server/http/router"
	"github.com/gennadyterekhov/metrics-storage/internal/server/repositories"
	"github.com/gennadyterekhov/metrics-storage/internal/server/services/services"
	"github.com/gennadyterekhov/metrics-storage/internal/server/storage"
)

// App is main instance of server app.
type App struct {
	Config      *config.ServerConfig
	DBOrRAM     storage.Interface
	Repository  repositories.RepositoryInterface
	Services    *services.Services
	Controllers *handlers.Controllers
	Router      *router.Router
}

// New creates App instance, injects all dependencies.
func New() *App {
	conf := config.New()
	DBOrRAM := storage.New(conf.DBDsn)
	repo := repositories.New(DBOrRAM)
	servicesPack := services.New(repo, conf)
	middlewareSet := middleware.New(conf)
	controllers := handlers.NewControllers(servicesPack, middlewareSet)
	rtr := router.New(controllers)

	return &App{
		Config:      conf,
		DBOrRAM:     DBOrRAM,
		Repository:  repo,
		Services:    servicesPack,
		Controllers: controllers,
		Router:      rtr,
	}
}

// StartServer starts a server; has graceful shutdown
func (a *App) StartServer(ctx context.Context) error {
	rootContext, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	a.tryToLoadFromDisk()

	if a.Config.StoreInterval != 0 {
		a.Services.TimeTracker.StartTrackingIntervals()
	}

	if a.Config.IsGrpc {
		srv := a.initGrpcServer()
		go a.gracefulShutdown(rootContext, nil, srv)
		return a.startGrpcServer(srv)
	}

	srv := a.initHTTPServer()
	go a.gracefulShutdown(rootContext, srv, nil)

	return a.startHTTPServer(srv)
}

func (a *App) tryToLoadFromDisk() {
	if a.Config.FileStorage != "" && a.Config.Restore {
		err := a.DBOrRAM.LoadFromDisk(context.Background(), a.Config.FileStorage)
		if err != nil {
			logger.Custom.Debugln("could not load metrics from disk, loaded empty repository")
		}
	}
}

// gracefulShutdown - this code runs if app gets any of (syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
func (a *App) gracefulShutdown(ctx context.Context, server *http.Server, grpcServer *grpc.Server) {
	<-ctx.Done()

	logger.Custom.Infoln("shutting down gracefully")

	a.Services.SaveMetricService.SaveToDisk(context.Background())

	if server != nil {
		err := server.Shutdown(ctx)
		if err != nil {
			logger.Custom.Errorln("error during server shutdown", err.Error())
		}
	}

	if grpcServer != nil {
		grpcServer.GracefulStop()
	}

	err := a.DBOrRAM.CloseDB()
	if err != nil {
		logger.Custom.Errorln("error during db shutdown", err.Error())
	}
}

func (a *App) initGrpcServer() *grpc.Server {
	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			loggerInterceptor.LoggingInterceptor,
			ipcontrol.New(a.Config.TrustedSubnet).IPControl,
		),
	)
	pb.RegisterMetricsServer(s, &grpcHandlers.Server{
		PingService:       a.Services.PingService,
		GetMetricService:  a.Services.GetMetricService,
		SaveMetricService: a.Services.SaveMetricService,
	})
	return s
}

func (a *App) initHTTPServer() *http.Server {
	server := &http.Server{}
	server.Handler = a.Router.ChiRouter
	server.Addr = a.Config.Addr
	return server
}

func (a *App) startGrpcServer(grpcServer *grpc.Server) error {
	listen, err := net.Listen("tcp", a.Config.Addr)
	if err != nil {
		return err
	}
	logger.Custom.Infoln("gRPC Server started on ", a.Config.Addr)

	return grpcServer.Serve(listen)
}

func (a *App) startHTTPServer(httpServer *http.Server) error {
	logger.Custom.Infoln("HTTP Server started on ", a.Config.Addr)
	return httpServer.ListenAndServe()
}
