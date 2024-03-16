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

	if storage.DBConnection == nil {
		http.Error(res, "db connection is not initialized", http.StatusInternalServerError)
	}

	err = storage.DBConnection.Ping()
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}

	res.WriteHeader(http.StatusOK)
	_, err = io.WriteString(res, "ok")
	if err != nil {
		logger.ZapSugarLogger.Errorln("error when writing ping response", err.Error())
	}
}
