package main

import (
	"fmt"
	"github.com/gennadyterekhov/metrics-storage/internal/server/httpui/handlers"
	"github.com/gennadyterekhov/metrics-storage/internal/server/storage"
	"net/http"
)

func main() {
	config := getConfig()
	fmt.Printf("Server started on %v\n", config.Addr)

	if config.FileStorage != "" {
		if config.Restore {
			err := storage.MetricsRepository.Load(config.FileStorage)
			if err != nil {
				panic(err)
				return
			}
		}
		defer storage.MetricsRepository.Save(config.FileStorage)
	}

	err := http.ListenAndServe(config.Addr, handlers.GetRouter())
	if err != nil {
		panic(err)
	}
}
