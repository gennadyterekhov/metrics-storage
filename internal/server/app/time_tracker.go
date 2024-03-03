package app

import (
	"github.com/gennadyterekhov/metrics-storage/internal/server/config"
	"github.com/gennadyterekhov/metrics-storage/internal/server/storage"
	"time"
)

type TimeTracker struct {
	ModOffset       int64 // this is timestamp % interval. when now % interval == ModOffset means interval has passed
	ActionFulfilled int
}

var TimeTrackerInstance = NewTimeTracker()

func StartTackingIntervals() {
	go routine()
}

func routine() {
	for range time.Tick(time.Duration(config.Conf.StoreInterval) * time.Second) {
		TimeTrackerInstance.onInterval()
	}
}

func NewTimeTracker() *TimeTracker {

	return &TimeTracker{
		ModOffset: time.Now().Unix() % int64(config.Conf.StoreInterval),
	}
}

// IsIntervalPassed is independent of run time  in contrary to time.Tick
func (ttr *TimeTracker) IsIntervalPassed() (ok bool) {

	return time.Now().Unix()%int64(config.Conf.StoreInterval) == ttr.ModOffset
}

func (ttr *TimeTracker) onInterval() {
	storage.MetricsRepository.Save(config.Conf.FileStorage)
	ttr.ActionFulfilled += 1
}
