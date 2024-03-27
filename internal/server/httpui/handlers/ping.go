package handlers

import (
	"github.com/gennadyterekhov/metrics-storage/internal/server/storage"
	_ "github.com/jackc/pgx/v5/stdlib"

	"net/http"
)

func Ping(res http.ResponseWriter, req *http.Request) {

	if storage.MetricsRepository == nil {
		http.Error(res, "storage is not initialized", http.StatusInternalServerError)
	}

	if !storage.MetricsRepository.IsDB() {
		http.Error(res, "storage is not of db type", http.StatusInternalServerError)
	}

}
