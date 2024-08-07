package app

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/gennadyterekhov/metrics-storage/internal/server/httpui/middleware"

	"github.com/gennadyterekhov/metrics-storage/internal/server/services/services"

	"github.com/gennadyterekhov/metrics-storage/internal/server/httpui/handlers/handlers"

	"github.com/gennadyterekhov/metrics-storage/internal/server/httpui/router"

	"github.com/gennadyterekhov/metrics-storage/internal/common/logger"
	"github.com/gennadyterekhov/metrics-storage/internal/server/config"
	"github.com/gennadyterekhov/metrics-storage/internal/server/repositories"
	"github.com/gennadyterekhov/metrics-storage/internal/server/storage"
)

// App is main instance of server app.
type App struct {
	Config      config.ServerConfig
	DBOrRAM     storage.Interface
	Repository  repositories.RepositoryInterface
	Services    services.Services
	Controllers handlers.Controllers
	Router      router.Router
}

// New creates App instance, injects all dependencies.
func New() *App {
	conf := config.New()
	DBOrRAM := storage.New(conf.DBDsn)
	repo := repositories.New(DBOrRAM)
	servicesPack := services.New(&repo, conf)
	middlewareSet := middleware.New(conf)
	controllers := handlers.NewControllers(&servicesPack, middlewareSet)
	rtr := router.New(&controllers)

	return &App{
		Config:      *conf,
		DBOrRAM:     DBOrRAM,
		Repository:  repo,
		Services:    servicesPack,
		Controllers: controllers,
		Router:      rtr,
	}
}

// StartServer starts a server; has graceful shutdown
func (a App) StartServer() error {
	var err error

	rootContext, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	if a.Config.FileStorage != "" {
		if a.Config.Restore {
			err = a.DBOrRAM.LoadFromDisk(context.Background(), a.Config.FileStorage)
			if err != nil {
				logger.Custom.Debugln("could not load metrics from disk, loaded empty repository")
			}
		}
	}

	if a.Config.StoreInterval != 0 {
		a.Services.TimeTracker.StartTrackingIntervals()
	}

	defer func(DBOrRAM storage.Interface) {
		errOnClose := DBOrRAM.CloseDB()
		if errOnClose != nil {
			logger.Custom.Errorln(errOnClose.Error())
		}
	}(a.DBOrRAM)

	_, err = fmt.Printf("Server started on %v\n", a.Config.Addr)
	if err != nil {
		return err
	}

	server := &http.Server{}
	server.Handler = a.Router.ChiRouter
	server.Addr = a.Config.Addr
	go a.gracefulShutdown(rootContext, server)

	err = server.ListenAndServe()
	if err != nil {
		return err
	}

	return err
}

// gracefulShutdown - this code runs if app gets any of (syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
func (a App) gracefulShutdown(ctx context.Context, server *http.Server) {
	<-ctx.Done()

	logger.Custom.Infoln("shutting down gracefully")

	a.Services.SaveMetricService.SaveToDisk(context.Background())
	err := server.Shutdown(ctx)
	if err != nil {
		logger.Custom.Errorln("error during server shutdown", err.Error())
	}
}
