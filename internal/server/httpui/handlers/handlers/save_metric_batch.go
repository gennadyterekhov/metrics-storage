package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gennadyterekhov/metrics-storage/internal/common/logger"
	"github.com/gennadyterekhov/metrics-storage/internal/server/httpui/middleware"
	"github.com/gennadyterekhov/metrics-storage/internal/server/httpui/requests"
)

func SaveMetricListHandler(cont SaveController) http.Handler {
	return middleware.CommonConveyor(
		http.HandlerFunc(cont.SaveMetricList),
	)
}

func (cont SaveController) SaveMetricList(res http.ResponseWriter, req *http.Request) {
	requestDto, err := getSaveListDtoForService(req)
	if err != nil {
		logger.ZapSugarLogger.Debugln("found error during request DTO build process", err.Error())
		writeErrorToOutput(&res, err)
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
