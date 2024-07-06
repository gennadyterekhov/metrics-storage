package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gennadyterekhov/metrics-storage/internal/server/services"

	"github.com/gennadyterekhov/metrics-storage/internal/common/constants"
	"github.com/gennadyterekhov/metrics-storage/internal/common/constants/exceptions"
	"github.com/gennadyterekhov/metrics-storage/internal/common/constants/types"
	"github.com/gennadyterekhov/metrics-storage/internal/common/logger"
	"github.com/gennadyterekhov/metrics-storage/internal/server/httpui/middleware"
	"github.com/gennadyterekhov/metrics-storage/internal/server/httpui/requests"
	"github.com/gennadyterekhov/metrics-storage/internal/server/httpui/responses"
	"github.com/gennadyterekhov/metrics-storage/internal/server/httpui/validators"
	"github.com/go-chi/chi/v5"
)

func SaveMetricHandler() http.Handler {
	return middleware.CommonConveyor(
		http.HandlerFunc(SaveMetric),
	)
}

func SaveMetric(res http.ResponseWriter, req *http.Request) {
	requestDto := getSaveDtoForService(req)
	if requestDto.Error != nil {
		logger.ZapSugarLogger.Debugln("found error during request DTO build process", requestDto.Error)
		writeErrorToOutput(&res, requestDto.Error)
		return
	}

	validatedRequestDto := validateSaveRequest(requestDto)
	if validatedRequestDto.Error != nil {
		logger.ZapSugarLogger.Debugln("found error during request validation", requestDto.Error)
		writeErrorToOutput(&res, validatedRequestDto.Error)
		return
	}

	responseDto := services.SaveMetricToMemory(req.Context(), requestDto)
	if responseDto.Error != nil {
		logger.ZapSugarLogger.Debugln(
			"found error during response DTO build process in usecase",
			requestDto.Error)
		writeErrorToOutput(&res, responseDto.Error)
		return
	}

	writeDtoToOutputIfJSON(&res, responseDto)
}

func getSaveDtoForService(req *http.Request) *requests.SaveMetricRequest {
	requestDto := &requests.SaveMetricRequest{
		IsJSON: false,
	}

	if req.Header.Get(constants.HeaderContentType) == constants.ApplicationJSON {
		requestDto.IsJSON = true
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(requestDto)
		requestDto.Error = err
		return requestDto
	}

	requestDto = validators.GetDataToSave(
		chi.URLParam(req, "metricType"),
		chi.URLParam(req, "metricName"),
		chi.URLParam(req, "metricValue"),
	)
	return requestDto
}

func validateSaveRequest(requestDto *requests.SaveMetricRequest) *requests.SaveMetricRequest {
	validatedRequestDto := requestDto
	if requestDto.MetricType != types.Counter && requestDto.MetricType != types.Gauge {
		validatedRequestDto.Error = fmt.Errorf(exceptions.InvalidMetricTypeChoice)
	}

	return validatedRequestDto
}

func writeDtoToOutputIfJSON(res *http.ResponseWriter, responseDto *responses.GetMetricResponse) {
	if responseDto.IsJSON {
		writeDtoToOutput(res, responseDto)
	}
}
