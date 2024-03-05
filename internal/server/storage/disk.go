package storage

import (
	"encoding/json"
	"github.com/gennadyterekhov/metrics-storage/internal/logger"
	"os"
)

func (strg *MemStorage) Save(filename string) (err error) {
	logger.ZapSugarLogger.Infoln("saving metrics to disk")

	data, err := json.MarshalIndent(strg, "", "   ")
	if err != nil {
		logger.ZapSugarLogger.Warnln("error when json encoding metrics")
		return err
	}

	err = os.WriteFile(filename, data, 0666)
	if err != nil {
		logger.ZapSugarLogger.Warnln("error when writing metrics file to disk")
		return err
	}
	return nil
}

func (strg *MemStorage) Load(filename string) (err error) {
	logger.ZapSugarLogger.Infoln("loading metrics from disk")

	fileBytes, err := os.ReadFile(filename)
	if err != nil {
		logger.ZapSugarLogger.Warnln("error when reading metrics file", err.Error())
		logger.ZapSugarLogger.Infoln("loading empty metrics")
		strg.Clear()
		strg = CreateStorage()

		return err
	}

	err = json.Unmarshal(fileBytes, strg)
	if err != nil {
		logger.ZapSugarLogger.Panicln("error when json decoding metrics from disk")
		return err
	}
	return nil
}
