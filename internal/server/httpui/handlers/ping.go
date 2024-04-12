package handlers

import (
	"github.com/gennadyterekhov/metrics-storage/internal/logger"
	"github.com/gennadyterekhov/metrics-storage/internal/server/storage"
	_ "github.com/jackc/pgx/v5/stdlib"

	"io"
	"net/http"
)

func Ping(res http.ResponseWriter, req *http.Request) {
	var err error

	if storage.MetricsRepository == nil {
		http.Error(res, "storage is not initialized", http.StatusInternalServerError)
	}

	dbStorage := storage.MetricsRepository.GetDB()

	if dbStorage != nil {
		err = dbStorage.DBConnection.Ping()
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}

		res.WriteHeader(http.StatusOK)
		_, err = io.WriteString(res, "ok")
		if err != nil {
			logger.ZapSugarLogger.Errorln("error when writing ping response", err.Error())
		}
	} else {
		http.Error(res, "DBStorage is nil: probably storage is not of db type", http.StatusInternalServerError)
	}
}
