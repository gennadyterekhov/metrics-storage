package handlers

import (
	"io"
	"net/http"

	"github.com/gennadyterekhov/metrics-storage/internal/common/constants"
	"github.com/gennadyterekhov/metrics-storage/internal/server/httpui/middleware"
	"github.com/gennadyterekhov/metrics-storage/internal/server/services/services"
	_ "github.com/gennadyterekhov/metrics-storage/swagger"
)

type GetController struct {
	Service       services.GetMetricService
	MiddlewareSet *middleware.Set
}

func NewGetController(serv services.GetMetricService, middlewareSet *middleware.Set) GetController {
	return GetController{
		Service:       serv,
		MiddlewareSet: middlewareSet,
	}
}

// GetAllMetrics returns html page with a list of all metrics with their values
// @Tags GET
// @Summary returns html page with a list of all metrics with their values
// @Description returns html page with a list of all metrics with their values
// @ID GetAllMetrics
// @Produce plain
// @Success 200 {object} string "ok"
// @Failure 500 {string} string "Internal server error"
// @Router / [get]
func (cont GetController) GetAllMetrics(res http.ResponseWriter, req *http.Request) {
	htmlPage := cont.Service.GetMetricsListAsHTML(req.Context())

	res.Header().Set(constants.HeaderContentType, constants.TextHTML)
	_, err := io.WriteString(res, htmlPage)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}

func GetAllMetricsHandler(cont GetController) http.Handler {
	return cont.MiddlewareSet.CommonConveyor(
		http.HandlerFunc(cont.GetAllMetrics),
	)
}
