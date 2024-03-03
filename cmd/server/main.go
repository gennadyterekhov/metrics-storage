package main

import (
	"fmt"
	"github.com/gennadyterekhov/metrics-storage/internal/logger"
	"github.com/gennadyterekhov/metrics-storage/internal/server/app"
	"github.com/gennadyterekhov/metrics-storage/internal/server/config"
	"github.com/gennadyterekhov/metrics-storage/internal/server/httpui/handlers"
	"github.com/gennadyterekhov/metrics-storage/internal/server/storage"
	"net/http"
)

func main() {
	fmt.Println("Starting")

	if config.Conf.FileStorage != "" {
		if config.Conf.Restore {
			err := storage.MetricsRepository.Load(config.Conf.FileStorage)
			if err != nil {
				logger.ZapSugarLogger.Debugln("could not load metrics from disk, but not panicking. just loaded empty repository")
				logger.ZapSugarLogger.Warnln("error when loading metrics from disk", err.Error())
			}
		}
		defer storage.MetricsRepository.Save(config.Conf.FileStorage)
	}

	if config.Conf.StoreInterval != 0 {
		app.StartTrackingIntervals()
	}

	fmt.Printf("Server started on %v\n", config.Conf.Addr)
	err := http.ListenAndServe(config.Conf.Addr, handlers.GetRouter())
	if err != nil {
		panic(err)
	}
}
