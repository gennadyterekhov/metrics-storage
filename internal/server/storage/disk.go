package storage

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/Rican7/retry"
	"github.com/Rican7/retry/backoff"
	"github.com/Rican7/retry/strategy"
	"github.com/gennadyterekhov/metrics-storage/internal/logger"
)

type SavedOnDisc struct {
	Counters map[string]int64   `json:"counters"`
	Gauges   map[string]float64 `json:"gauges"`
}

func (strg *MemStorage) SaveToDisk(ctx context.Context, filename string) (err error) {
	logger.ZapSugarLogger.Infoln("saving metrics to disk")

	data, err := json.MarshalIndent(strg, "", "   ")
	if err != nil {
		logger.ZapSugarLogger.Errorln("error when json encoding metrics")
		return err
	}

	err = os.WriteFile(filename, data, 0o666)
	if err != nil {
		logger.ZapSugarLogger.Errorln("error when writing metrics file to disk")
		return err
	}
	return nil
}

func (strg *MemStorage) LoadFromDisk(ctx context.Context, filename string) (err error) {
	logger.ZapSugarLogger.Infoln("loading metrics from disk")

	err = strg.loadFromDiskWithRetry(filename)
	if err != nil {
		logger.ZapSugarLogger.Errorln("error when reading metrics file", err.Error())
		logger.ZapSugarLogger.Infoln("loading empty metrics")
		strg.Clear()
		MetricsRepository = CreateStorage()

		return err
	}

	return nil
}

func (strg *MemStorage) loadFromDiskWithRetry(filename string) error {
	return retry.Retry(
		func(attempt uint) error {
			logger.ZapSugarLogger.Debugf("loading metrics from disk attempt: %v", attempt)

			fileBytes, err := os.ReadFile(filename)
			if err != nil {
				logger.ZapSugarLogger.Errorln("error when reading metrics file", err.Error())
				return err
			}

			err = json.Unmarshal(fileBytes, strg)
			if err != nil {
				logger.ZapSugarLogger.Errorln("error when json decoding metrics from disk.", err.Error())
				return err
			}
			return nil
		},
		strategy.Limit(3),
		strategy.Backoff(backoff.Incremental(0*time.Second, 3*time.Second)),
	)
}

func (strg *DBStorage) SaveToDisk(ctx context.Context, filename string) (err error) {
	logger.ZapSugarLogger.Infoln("saving metrics to disk")

	savedOnDisc := &SavedOnDisc{}
	savedOnDisc.Gauges = strg.GetAllGauges(ctx)
	savedOnDisc.Counters = strg.GetAllCounters(ctx)

	data, err := json.MarshalIndent(savedOnDisc, "", "   ")
	if err != nil {
		logger.ZapSugarLogger.Warnln("error when json encoding metrics")
		return err
	}

	err = os.WriteFile(filename, data, 0o666)
	if err != nil {
		logger.ZapSugarLogger.Warnln("error when writing metrics file to disk")
		return err
	}
	return nil
}

func (strg *DBStorage) LoadFromDisk(ctx context.Context, filename string) (err error) {
	logger.ZapSugarLogger.Infoln("will not load from disk to db, database is already persistent storage")

	return nil
}
