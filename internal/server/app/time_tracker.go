package app

import (
	"github.com/gennadyterekhov/metrics-storage/internal/logger"
	"github.com/gennadyterekhov/metrics-storage/internal/server/config"
	"github.com/gennadyterekhov/metrics-storage/internal/server/storage"
	"time"
)

type TimeTracker struct {
	ModOffset       int64 // this is timestamp % interval. when now % interval == ModOffset means interval has passed
	ActionFulfilled int
}

var TimeTrackerInstance = NewTimeTracker()

func StartTrackingIntervals() {
	go routine()
}

func routine() {
	if config.Conf.StoreInterval == 0 {
		return
	}
	for range time.Tick(time.Duration(config.Conf.StoreInterval) * time.Second) {
		TimeTrackerInstance.onInterval()
	}
}

func NewTimeTracker() *TimeTracker {
	var offset int64 = 0
	if config.Conf.StoreInterval != 0 {
		offset = time.Now().Unix() % int64(config.Conf.StoreInterval)
	}
	return &TimeTracker{
		ModOffset: offset,
	}
}

// IsIntervalPassed is independent of run time  in contrary to time.Tick
func (ttr *TimeTracker) IsIntervalPassed() (ok bool) {

	return time.Now().Unix()%int64(config.Conf.StoreInterval) == ttr.ModOffset
}

func (ttr *TimeTracker) onInterval() {
	logger.ZapSugarLogger.Infoln("STORE_INTERVAL passed, saving metrics to disk")
	storage.MetricsRepository.Save(config.Conf.FileStorage)
	ttr.ActionFulfilled += 1
}
