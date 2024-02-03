package main

import (
	"github.com/gennadyterekhov/metrics-storage/internal/handlers"
	"net/http"
)

func registerHandlers(mux *http.ServeMux) {
	mux.HandleFunc(`/update/`, handlers.SaveMetric)
}

func main() {
	mux := http.NewServeMux()
	registerHandlers(mux)

	err := http.ListenAndServe(`0.0.0.0:8080`, mux)
	if err != nil {
		panic(err)
	}
}
