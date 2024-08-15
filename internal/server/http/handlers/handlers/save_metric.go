package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gennadyterekhov/metrics-storage/internal/server/http/requests"

	"github.com/gennadyterekhov/metrics-storage/internal/common/constants"
	"github.com/gennadyterekhov/metrics-storage/internal/common/constants/exceptions"
	"github.com/gennadyterekhov/metrics-storage/internal/common/constants/types"
	"github.com/gennadyterekhov/metrics-storage/internal/common/logger"
	"github.com/gennadyterekhov/metrics-storage/internal/server/http/middleware"
	"github.com/gennadyterekhov/metrics-storage/internal/server/http/responses"
	"github.com/gennadyterekhov/metrics-storage/internal/server/http/validators"
	"github.com/gennadyterekhov/metrics-storage/internal/server/services/services"
	_ "github.com/gennadyterekhov/metrics-storage/swagger" // swagger
	"github.com/go-chi/chi/v5"
)

type SaveController struct {
	Service       *services.SaveMetricService
	MiddlewareSet *middleware.Set
}

func NewSaveController(serv *services.SaveMetricService, middlewareSet *middleware.Set) *SaveController {
	return &SaveController{
		Service:       serv,
		MiddlewareSet: middlewareSet,
	}
}

func SaveMetricHandler(cont *SaveController) http.Handler {
	return cont.MiddlewareSet.CommonConveyor(
		http.HandlerFunc(cont.SaveMetric),
	)
}

func SaveMetricJSONHandler(cont *SaveController) http.Handler {
	return cont.MiddlewareSet.CommonConveyor(
		http.HandlerFunc(cont.SaveMetricJSON),
	)
}

// SaveMetric saves metric to db. returns json with saved metric
// @Tags POST
// @Summary saves metric to db
// @Description saves metric to db
// @ID SaveMetric
// @Accept  plain
// @Produce plain
// @Param metricType path string true "'gauge' or 'counter'"
// @Param metricName path string true "name of metric, serves as identifier"
// @Param metricValue path string true "int64 if type is counter, float64 if type is gauge"
// @Success 200 {object} string "ok"
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "unknown metric type"
// @Failure 500 {string} string "Internal server error"
// @Router /update/{metricType}/{metricName}/{metricValue} [post]
func (cont SaveController) SaveMetric(res http.ResponseWriter, req *http.Request) {
	cont.saveMetricCommon(res, req)
}

// SaveMetricJSON saves metric to db. returns json with saved metric
// @Tags POST
// @Summary saves metric to db
// @Description saves metric to db
// @ID SaveMetricJSON
// @Accept  json
// @Produce json
// @Param   data body string true "requests.SaveMetricRequest"
// @Success 200 {object} string "ok"
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "unknown metric type"
// @Failure 500 {string} string "Internal server error"
// @Router /update [post]
func (cont SaveController) SaveMetricJSON(res http.ResponseWriter, req *http.Request) {
	cont.saveMetricCommon(res, req)
}

func (cont SaveController) saveMetricCommon(res http.ResponseWriter, req *http.Request) {
	requestDto := cont.getSaveDtoForService(req)
	if requestDto.Error != nil {
		logger.Custom.Debugln("found error during request DTO build process", requestDto.Error)
		writeErrorToOutput(res, requestDto.Error)
		return
	}

	validatedRequestDto := cont.validateSaveRequest(requestDto)
	if validatedRequestDto.Error != nil {
		logger.Custom.Debugln("found error during request validation", requestDto.Error)
		writeErrorToOutput(res, validatedRequestDto.Error)
		return
	}

	responseDto := cont.Service.SaveMetricToMemory(req.Context(), requestDto)
	if responseDto.Error != nil {
		logger.Custom.Debugln(
			"found error during response DTO build process in usecase",
			requestDto.Error)
		writeErrorToOutput(res, responseDto.Error)
		return
	}

	cont.writeDtoToOutputIfJSON(res, responseDto)
}

func (cont SaveController) getSaveDtoForService(req *http.Request) *requests.SaveMetricRequest {
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

func (cont SaveController) validateSaveRequest(requestDto *requests.SaveMetricRequest) *requests.SaveMetricRequest {
	validatedRequestDto := requestDto
	if requestDto.MetricType != types.Counter && requestDto.MetricType != types.Gauge {
		validatedRequestDto.Error = fmt.Errorf(exceptions.InvalidMetricTypeChoice)
	}

	return validatedRequestDto
}

func (cont SaveController) writeDtoToOutputIfJSON(res http.ResponseWriter, responseDto *responses.GetMetricResponse) {
	if responseDto.IsJSON {
		writeDtoToOutput(res, responseDto)
	}
}
