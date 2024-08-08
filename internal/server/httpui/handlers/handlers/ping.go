package handlers

import (
	"io"
	"net/http"

	"github.com/gennadyterekhov/metrics-storage/internal/server/services/services"

	"github.com/gennadyterekhov/metrics-storage/internal/common/logger"
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
func (cont *PingController) Ping(res http.ResponseWriter, _ *http.Request) {
	var err error

	if cont.Service.Repository == nil {
		http.Error(res, "storage is not initialized", http.StatusInternalServerError)
	}
	res.WriteHeader(http.StatusOK)
	_, err = io.WriteString(res, "ok")
	if err != nil {
		logger.Custom.Errorln("error when writing ping response", err.Error())
	}
}
