package storage

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/Rican7/retry"
	"github.com/Rican7/retry/backoff"
	"github.com/Rican7/retry/strategy"
	"github.com/gennadyterekhov/metrics-storage/internal/common/logger"
)

type SavedOnDisc struct {
	Counters map[string]int64   `json:"counters"`
	Gauges   map[string]float64 `json:"gauges"`
}

func (strg *MemStorage) SaveToDisk(_ context.Context, filename string) (err error) {
	logger.Custom.Infoln("saving metrics to disk")

	data, err := json.MarshalIndent(strg, "", "   ")
	if err != nil {
		logger.Custom.Errorln("error when json encoding metrics")
		return err
	}

	err = os.WriteFile(filename, data, 0o666)
	if err != nil {
		logger.Custom.Errorln("error when writing metrics file to disk")
		return err
	}
	return nil
}

func (strg *MemStorage) LoadFromDisk(_ context.Context, filename string) (err error) {
	logger.Custom.Infoln("loading metrics from disk")

	err = strg.loadFromDiskWithRetry(filename)
	if err != nil {
		logger.Custom.Debugln("error when reading metrics file", err.Error())
		logger.Custom.Infoln("loading empty metrics")
		strg.Clear()

		return err
	}

	return nil
}

func (strg *MemStorage) loadFromDiskWithRetry(filename string) error {
	return retry.Retry(
		func(attempt uint) error {
			logger.Custom.Debugf("loading metrics from disk attempt: %v", attempt)

			fileBytes, err := os.ReadFile(filename)
			if err != nil {
				logger.Custom.Debugln("error when reading metrics file", err.Error())
				return err
			}

			err = json.Unmarshal(fileBytes, strg)
			if err != nil {
				logger.Custom.Debugln("error when json decoding metrics from disk.", err.Error())
				return err
			}
			return nil
		},
		strategy.Limit(3),
		strategy.Backoff(backoff.Incremental(0*time.Second, 3*time.Second)),
	)
}

func (strg *DBStorage) SaveToDisk(ctx context.Context, filename string) (err error) {
	logger.Custom.Infoln("saving metrics to disk")

	savedOnDisc := &SavedOnDisc{}
	savedOnDisc.Gauges = strg.GetAllGauges(ctx)
	savedOnDisc.Counters = strg.GetAllCounters(ctx)

	data, err := json.MarshalIndent(savedOnDisc, "", "   ")
	if err != nil {
		logger.Custom.Errorln("error when json encoding metrics")
		return err
	}

	err = os.WriteFile(filename, data, 0o666)
	if err != nil {
		logger.Custom.Errorln("error when writing metrics file to disk")
		return err
	}
	return nil
}

func (strg *DBStorage) LoadFromDisk(_ context.Context, _ string) (err error) {
	logger.Custom.Infoln("will not load from disk to db, database is already persistent storage")

	return nil
}
