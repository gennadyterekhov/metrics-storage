package services

import (
	"context"
	"time"

	"github.com/gennadyterekhov/metrics-storage/internal/common/logger"
	"github.com/gennadyterekhov/metrics-storage/internal/server/config"
	"github.com/gennadyterekhov/metrics-storage/internal/server/storage"
)

func StartTrackingIntervals() {
	ticker := time.NewTicker(time.Duration(config.Conf.StoreInterval) * time.Second)
	go routine(ticker)
}

func routine(ticker *time.Ticker) {
	if config.Conf.StoreInterval == 0 {
		return
	}
	for {
		<-ticker.C
		onInterval()
	}
}

func onInterval() {
	logger.ZapSugarLogger.Infoln("STORE_INTERVAL passed, saving metrics to disk")
	storage.MetricsRepository.SaveToDisk(context.Background(), config.Conf.FileStorage)
}
