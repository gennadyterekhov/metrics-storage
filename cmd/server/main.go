package main

import (
	"fmt"
	"github.com/gennadyterekhov/metrics-storage/internal/logger"
	"github.com/gennadyterekhov/metrics-storage/internal/server/app"
	"github.com/gennadyterekhov/metrics-storage/internal/server/config"
	"github.com/gennadyterekhov/metrics-storage/internal/server/httpui/handlers"
	"github.com/gennadyterekhov/metrics-storage/internal/server/storage"
	"net/http"
	"os"
	"os/signal"
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
	}

	if config.Conf.StoreInterval != 0 {
		app.StartTrackingIntervals()
	}

	go onStop()
	fmt.Printf("Server started on %v\n", config.Conf.Addr)
	err := http.ListenAndServe(config.Conf.Addr, handlers.GetRouter())

	if err != nil {
		panic(err)
	}
}

func onStop() {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)
	<-sigchan
	logger.ZapSugarLogger.Infoln("shutting down gracefully")

	err := storage.MetricsRepository.Save(config.Conf.FileStorage)
	if err != nil {
		logger.ZapSugarLogger.Errorln("could not save metrics during shutdown")
	}
	os.Exit(0)
}
