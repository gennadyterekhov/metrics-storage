package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

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
func New() App {
	conf := config.New()
	DBOrRAM := storage.New(conf.DBDsn)
	repo := repositories.New(DBOrRAM)
	servicesPack := services.New(&repo, &conf)
	middlewareSet := middleware.New(&conf)
	controllers := handlers.NewControllers(&servicesPack, middlewareSet)
	rtr := router.New(&controllers)

	return App{
		Config:      conf,
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

	if a.Config.FileStorage != "" {
		if a.Config.Restore {
			err = a.DBOrRAM.LoadFromDisk(context.Background(), a.Config.FileStorage)
			if err != nil {
				logger.ZapSugarLogger.Debugln("could not load metrics from disk, loaded empty repository")
				logger.ZapSugarLogger.Errorln("error when loading metrics from disk", err.Error())
			}
		}
	}

	if a.Config.StoreInterval != 0 {
		a.Services.TimeTracker.StartTrackingIntervals()
	}

	defer func(DBOrRAM storage.Interface) {
		err := DBOrRAM.CloseDB()
		if err != nil {
			fmt.Println(err.Error())
		}
	}(a.DBOrRAM)

	go a.onStop()
	fmt.Printf("Server started on %v\n", a.Config.Addr)
	err = http.ListenAndServe(a.Config.Addr, a.Router.ChiRouter)

	return err
}

func (a App) onStop() {
	sigchan := make(chan os.Signal, 1)
	defer close(sigchan)
	signal.Notify(sigchan, os.Interrupt)
	<-sigchan
	logger.ZapSugarLogger.Infoln("shutting down gracefully")

	a.Services.SaveMetricService.SaveToDisk(context.Background())
}
