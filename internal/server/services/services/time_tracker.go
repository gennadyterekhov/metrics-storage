package services

import (
	"context"
	"time"

	"github.com/gennadyterekhov/metrics-storage/internal/server/repositories"

	"github.com/gennadyterekhov/metrics-storage/internal/common/logger"
	"github.com/gennadyterekhov/metrics-storage/internal/server/config"
)

type TimeTracker struct {
	Repository repositories.RepositoryInterface
	Config     *config.ServerConfig
}

func NewTimeTracker(repo repositories.RepositoryInterface, conf *config.ServerConfig) TimeTracker {
	return TimeTracker{
		Repository: repo,
		Config:     conf,
	}
}

func (tt TimeTracker) StartTrackingIntervals() {
	ticker := time.NewTicker(time.Duration(tt.Config.StoreInterval) * time.Second)
	go tt.routine(ticker)
}

func (tt TimeTracker) routine(ticker *time.Ticker) {
	if tt.Config.StoreInterval == 0 {
		return
	}
	for {
		<-ticker.C
		tt.onInterval()
	}
}

func (tt TimeTracker) onInterval() {
	logger.ZapSugarLogger.Infoln("STORE_INTERVAL passed, saving metrics to disk")
	err := tt.Repository.SaveToDisk(context.Background(), tt.Config.FileStorage)
	if err != nil {
		logger.ZapSugarLogger.Errorln("error when saving to disk on interval")
	}
}
