package handlers

import (
	"github.com/go-chi/chi/v5"

	"github.com/gennadyterekhov/metrics-storage/internal/server/config"
	"github.com/gennadyterekhov/metrics-storage/internal/server/httpui/middleware"
	"github.com/gennadyterekhov/metrics-storage/internal/server/repositories"
	"github.com/gennadyterekhov/metrics-storage/internal/server/services/services"
	"github.com/gennadyterekhov/metrics-storage/internal/server/storage"
)

func SetUpExampleRouter() (*chi.Mux, Controllers) {
	serverConfig := config.New()
	repo := repositories.New(storage.New(""))
	servs := services.New(repo, serverConfig)
	middlewareSet := middleware.New(serverConfig)
	controllersStruct := NewControllers(&servs, middlewareSet)
	chiRouter := chi.NewRouter()
	return chiRouter, controllersStruct
}
