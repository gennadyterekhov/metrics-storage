package validators

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/gennadyterekhov/metrics-storage/internal/common/constants/exceptions"
	"github.com/gennadyterekhov/metrics-storage/internal/common/constants/types"
	"github.com/gennadyterekhov/metrics-storage/internal/common/logger"
	"github.com/gennadyterekhov/metrics-storage/internal/server/httpui/requests"
)

func GetDataToSave(metricType string, metricName string, metricValue string) *requests.SaveMetricRequest {
	filledDto := &requests.SaveMetricRequest{}
	if metricType == "" {
		filledDto.Error = fmt.Errorf(exceptions.EmptyMetricType)
	}
	if metricName == "" {
		filledDto.Error = fmt.Errorf(exceptions.EmptyMetricName)
	}
	if metricValue == "" {
		filledDto.Error = fmt.Errorf(exceptions.EmptyMetricValue)
	}
	counterValue, gaugeValue, err := validateParameters(metricType, metricName, metricValue)
	if err != nil {
		filledDto.Error = fmt.Errorf(exceptions.InvalidMetricValueFormat)
		logger.ZapSugarLogger.Errorln("error when building save request dto", filledDto.Error, err)
	}
	filledDto.MetricName = metricName
	filledDto.MetricType = metricType
	filledDto.CounterValue = &counterValue
	filledDto.GaugeValue = &gaugeValue

	return filledDto
}

func GetDataToGet(metricType string, metricName string) (typ string, name string, err error) {
	if metricType == "" {
		return "", "", fmt.Errorf(exceptions.EmptyMetricType)
	}
	if metricName == "" {
		return "", "", fmt.Errorf(exceptions.EmptyMetricName)
	}

	err = validateMetricType(metricType)
	if err != nil {
		return "", "", err
	}

	return metricType, metricName, nil
}

func validateParameters(metricTypeRaw string, nameRaw string, valueRaw string) (int64, float64, error) {
	err := validateMetricType(metricTypeRaw)
	if err != nil {
		return 0, 0, err
	}
	err = validateMetricName(nameRaw)
	if err != nil {
		return 0, 0, err
	}
	counterValue, gaugeValue, err := validateMetricValue(metricTypeRaw, valueRaw)
	if err != nil {
		return 0, 0, err
	}

	return counterValue, gaugeValue, nil
}

func validateMetricType(metricTypeRaw string) (err error) {
	if metricTypeRaw != types.Counter && metricTypeRaw != types.Gauge {
		return fmt.Errorf(exceptions.InvalidMetricTypeChoice)
	}
	return nil
}

func validateMetricName(nameRaw string) (err error) {
	if len(nameRaw) < 1 {
		return errors.New(exceptions.EmptyMetricName)
	}
	return nil
}

func validateMetricValue(metricTypeValidated string, valueRaw string) (counterValue int64, gaugeValue float64, err error) {
	if metricTypeValidated == types.Counter {
		val, err := strconv.ParseInt(valueRaw, 10, 64)
		if err != nil {
			return 0, 0, err
		}
		return val, 0, nil
	}
	if metricTypeValidated == types.Gauge {
		val, err := strconv.ParseFloat(valueRaw, 64)
		if err != nil {
			return 0, 0, err
		}
		return 0, val, nil
	}
	return 0, 0, errors.New("unexpected type after validation")
}
