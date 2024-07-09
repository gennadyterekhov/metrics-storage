package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gennadyterekhov/metrics-storage/internal/common/constants"
	"github.com/gennadyterekhov/metrics-storage/internal/common/constants/exceptions"
	"github.com/gennadyterekhov/metrics-storage/internal/common/constants/types"
	"github.com/gennadyterekhov/metrics-storage/internal/common/logger"
	"github.com/gennadyterekhov/metrics-storage/internal/server/httpui/requests"
	"github.com/gennadyterekhov/metrics-storage/internal/server/httpui/responses"
	"github.com/gennadyterekhov/metrics-storage/internal/server/httpui/validators"
	"github.com/go-chi/chi/v5"
)

func GetMetricHandler(cont GetController) http.Handler {
	return cont.MiddlewareSet.CommonConveyor(
		http.HandlerFunc(cont.GetMetric),
	)
}

func (cont GetController) GetMetric(res http.ResponseWriter, req *http.Request) {
	requestDto := cont.getDtoForService(req)
	if requestDto.Error != nil {
		logger.ZapSugarLogger.Debugln("found error during request DTO build process", requestDto.Error)
		writeErrorToOutput(&res, requestDto.Error)
		return
	}

	validatedRequestDto := cont.validateRequest(requestDto)
	if validatedRequestDto.Error != nil {
		logger.ZapSugarLogger.Debugln("found error during request validation", validatedRequestDto.Error)
		writeErrorToOutput(&res, validatedRequestDto.Error)
		return
	}

	responseDto := cont.Service.GetMetric(req.Context(), requestDto)
	if responseDto.Error != nil {
		logger.ZapSugarLogger.Debugln("found error during response DTO build process in usecase", responseDto.Error)
		writeErrorToOutput(&res, responseDto.Error)
		return
	}

	writeDtoToOutput(&res, responseDto)
}

func (cont GetController) getDtoForService(req *http.Request) *requests.GetMetricRequest {
	dto := &requests.GetMetricRequest{
		IsJSON: false,
	}

	if req.Header.Get(constants.HeaderContentType) == constants.ApplicationJSON {
		dto.IsJSON = true
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(dto)
		dto.Error = err
		return dto
	}

	metricType, name, err := validators.GetDataToGet(
		chi.URLParam(req, "metricType"),
		chi.URLParam(req, "metricName"),
	)
	dto.Error = err
	dto.MetricName = name
	dto.MetricType = metricType

	return dto
}

func writeDtoToOutput(res *http.ResponseWriter, responseDto *responses.GetMetricResponse) {
	if responseDto.IsJSON {
		(*res).Header().Set(constants.HeaderContentType, constants.ApplicationJSON)
	}

	responseBody := serializeDto(responseDto)

	logger.ZapSugarLogger.Infoln("successfully serialized response body", string(responseBody))

	var err error
	_, err = io.WriteString(*res, string(responseBody))
	if err != nil {
		http.Error(*res, err.Error(), http.StatusInternalServerError)
		return
	}
	logger.ZapSugarLogger.Infoln("successfully written response body")
}

func serializeDto(responseDto *responses.GetMetricResponse) []byte {
	if responseDto.IsJSON {
		responseBytes, err := json.Marshal(responseDto)
		if err != nil {
			logger.ZapSugarLogger.Errorln("error when encoding json response body", err.Error())

			return []byte(err.Error())
		}

		return responseBytes
	}

	if responseDto.MetricType == types.Counter {
		return []byte(strconv.FormatInt(*responseDto.CounterValue, 10))
	}
	return []byte(strconv.FormatFloat(*responseDto.GaugeValue, 'g', -1, 64))
}

func writeErrorToOutput(res *http.ResponseWriter, err error) {
	logger.ZapSugarLogger.Debugln("writing error to output", err.Error())

	code := http.StatusInternalServerError
	if err.Error() == exceptions.UnknownMetricName {
		code = http.StatusNotFound
	}
	if err.Error() == exceptions.UnknownMetricType {
		code = http.StatusNotFound
	}
	if err.Error() == exceptions.InvalidMetricTypeChoice {
		code = http.StatusNotFound
	}
	if err.Error() == exceptions.InvalidMetricType {
		code = http.StatusNotFound
	}
	if err.Error() == exceptions.InvalidMetricValueFormat {
		code = http.StatusBadRequest
	}
	logger.ZapSugarLogger.Debugln("selected http error code ", code)

	http.Error(*res, err.Error(), code)
}

func (cont GetController) validateRequest(requestDto *requests.GetMetricRequest) *requests.GetMetricRequest {
	validatedRequestDto := requestDto
	if requestDto.MetricType != types.Counter && requestDto.MetricType != types.Gauge {
		validatedRequestDto.Error = fmt.Errorf(exceptions.InvalidMetricTypeChoice)
	}

	return validatedRequestDto
}
