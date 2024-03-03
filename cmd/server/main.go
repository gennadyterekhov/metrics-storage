package main

import (
	"fmt"
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
				panic(err)
				return
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
