package handlers

import (
	"io"
	"net/http"

	"github.com/gennadyterekhov/metrics-storage/internal/server/services/services"

	"github.com/gennadyterekhov/metrics-storage/internal/common/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func PingHandler(cont PingController) http.Handler {
	return http.HandlerFunc(cont.Ping)
}

type PingController struct {
	Service services.PingService
}

func NewPingController(serv services.PingService) PingController {
	return PingController{
		Service: serv,
	}
}

// Ping check db connection
// @Summary check db connection
// @Description check db connection
// @ID Ping
// @Accept  plain
// @Produce plain
// @Success 200 {object} string "ok"
// @Failure 500 {string} string "Internal server error"
// @Router /ping [get]
func (cont *PingController) Ping(res http.ResponseWriter, req *http.Request) {
	var err error

	if cont.Service.Repository == nil {
		http.Error(res, "storage is not initialized", http.StatusInternalServerError)
	}
	res.WriteHeader(http.StatusOK)
	_, err = io.WriteString(res, "ok")
	if err != nil {
		logger.ZapSugarLogger.Errorln("error when writing ping response", err.Error())
	}
	//dbStorage := cont.Service.Repository.GetDB()
	//
	//if dbStorage != nil {
	//	err = dbStorage.DBConnection.Ping()
	//	if err != nil {
	//		http.Error(res, err.Error(), http.StatusInternalServerError)
	//	}
	//
	//	res.WriteHeader(http.StatusOK)
	//	_, err = io.WriteString(res, "ok")
	//	if err != nil {
	//		logger.ZapSugarLogger.Errorln("error when writing ping response", err.Error())
	//	}
	//} else {
	//	http.Error(res, "DBStorage is nil: probably storage is not of db type", http.StatusInternalServerError)
	//}
}
