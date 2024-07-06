package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/gennadyterekhov/metrics-storage/internal/server/httpui/router"

	"github.com/gennadyterekhov/metrics-storage/internal/common/logger"
	"github.com/gennadyterekhov/metrics-storage/internal/server/config"
	"github.com/gennadyterekhov/metrics-storage/internal/server/httpui/handlers"
	"github.com/gennadyterekhov/metrics-storage/internal/server/repositories"
	"github.com/gennadyterekhov/metrics-storage/internal/server/services"
	"github.com/gennadyterekhov/metrics-storage/internal/server/storage"
)

type App struct {
	Config      config.ServerConfig
	DBOrRam     storage.StorageInterface
	Repository  repositories.RepositoryInterface
	Services    services.Services
	Controllers handlers.Controllers
	Router      router.Router
}

func New() App {
	conf := config.New()
	DBOrRam := storage.New(&conf)
	repo := repositories.New(DBOrRam)
	servicesPack := services.NewServices(&repo, &conf)
	controllers := handlers.NewControllers(&servicesPack)
	rtr := router.New(&controllers)

	return App{
		Config:      conf,
		DBOrRam:     DBOrRam,
		Repository:  repo,
		Services:    servicesPack,
		Controllers: controllers,
		Router:      rtr,
	}
}

func (a App) StartServer() error {
	var err error

	if a.Config.FileStorage != "" {
		if a.Config.Restore {
			err = a.DBOrRam.LoadFromDisk(context.Background(), a.Config.FileStorage)
			if err != nil {
				logger.ZapSugarLogger.Debugln("could not load metrics from disk, loaded empty repository")
				logger.ZapSugarLogger.Errorln("error when loading metrics from disk", err.Error())
			}
		}
	}

	if a.Config.StoreInterval != 0 {
		a.Services.TimeTracker.StartTrackingIntervals()
	}

	defer a.DBOrRam.CloseDB()

	go a.onStop()
	fmt.Printf("Server started on %v\n", a.Config.Addr)
	err = http.ListenAndServe(a.Config.Addr, a.Router.ChiRouter)

	return nil
}

func (a App) onStop() {
	sigchan := make(chan os.Signal, 1)
	defer close(sigchan)
	signal.Notify(sigchan, os.Interrupt)
	<-sigchan
	logger.ZapSugarLogger.Infoln("shutting down gracefully")

	a.Services.SaveMetricService.SaveToDisk(context.Background())
	os.Exit(0)
}
