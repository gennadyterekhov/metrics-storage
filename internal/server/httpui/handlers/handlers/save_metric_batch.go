package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gennadyterekhov/metrics-storage/internal/common/logger"
	"github.com/gennadyterekhov/metrics-storage/internal/server/httpui/requests"
)

func SaveMetricListHandler(cont SaveController) http.Handler {
	return cont.MiddlewareSet.CommonConveyor(
		http.HandlerFunc(cont.SaveMetricList),
	)
}

// SaveMetricList saves metric batch to db.
// @Tags POST
// @Summary saves metric batch to db
// @Description saves metric batch to db
// @ID SaveMetricList
// @Accept  json
// @Produce json
// @Param   data body string true "requests.SaveMetricListRequest"
// @Success 200 {object} string "ok"
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "unknown metric type"
// @Failure 500 {string} string "Internal server error"
// @Router /updates [post]
func (cont SaveController) SaveMetricList(res http.ResponseWriter, req *http.Request) {
	requestDto, err := getSaveListDtoForService(req)
	if err != nil {
		logger.Custom.Debugln("found error during request DTO build process", err.Error())
		writeErrorToOutput(res, err)
		return
	}

	cont.Service.SaveMetricListToMemory(req.Context(), requestDto)
	res.WriteHeader(http.StatusOK)
}

func getSaveListDtoForService(req *http.Request) (*requests.SaveMetricListRequest, error) {
	requestDto := &requests.SaveMetricListRequest{}

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(requestDto)
	if err != nil {
		return nil, err
	}

	return requestDto, nil
}
