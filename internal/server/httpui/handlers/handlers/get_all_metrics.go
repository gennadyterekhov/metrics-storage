package handlers

import (
	"io"
	"net/http"

	"github.com/gennadyterekhov/metrics-storage/internal/server/services/services"

	"github.com/gennadyterekhov/metrics-storage/internal/common/constants"
	"github.com/gennadyterekhov/metrics-storage/internal/server/httpui/middleware"
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
