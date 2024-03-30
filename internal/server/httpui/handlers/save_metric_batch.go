package handlers

import (
	"encoding/json"
	"github.com/gennadyterekhov/metrics-storage/internal/logger"
	"github.com/gennadyterekhov/metrics-storage/internal/server/app"
	"github.com/gennadyterekhov/metrics-storage/internal/server/httpui/middleware"
	"github.com/gennadyterekhov/metrics-storage/internal/server/httpui/requests"
	"net/http"
)

func SaveMetricListHandler() http.Handler {
	return middleware.CommonConveyor(
		http.HandlerFunc(SaveMetricList),
	)
}

func SaveMetricList(res http.ResponseWriter, req *http.Request) {
	requestDto, err := getSaveListDtoForService(req)
	if err != nil {
		logger.ZapSugarLogger.Debugln("found error during request DTO build process", err.Error())
		writeErrorToOutput(&res, err)
		return
	}

	app.SaveMetricListToMemory(requestDto)
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
