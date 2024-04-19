package agent

import (
	"context"
	"github.com/gennadyterekhov/metrics-storage/internal/agent/client"
	"github.com/gennadyterekhov/metrics-storage/internal/agent/metric"
	"github.com/gennadyterekhov/metrics-storage/internal/agent/poller"
	"github.com/gennadyterekhov/metrics-storage/internal/agent/sender"
	"github.com/gennadyterekhov/metrics-storage/internal/logger"
	"time"
)

type AgentConfig struct {
	Addr                      string
	IsGzip                    bool
	ReportInterval            int
	PollInterval              int
	IsBatch                   bool
	PayloadSignatureKey       string
	SimultaneousRequestsLimit int
}

func RunAgent(ctx context.Context, config *AgentConfig) {
	metricsSet := &metric.MetricsSet{}

	pollerInstance := poller.PollMaker{
		MetricsSet: metricsSet,
		Interval:   config.PollInterval,
		IsRunning:  false,
	}
	senderInstance := sender.MetricsSender{
		Address:         config.Addr,
		IsGzip:          config.IsGzip,
		Interval:        config.ReportInterval,
		IsRunning:       false,
		IsBatch:         config.IsBatch,
		NumberOfWorkers: config.SimultaneousRequestsLimit,
	}
	metricsStorageClient := client.MetricsStorageClient{
		Address:             config.Addr,
		IsGzip:              config.IsGzip,
		PayloadSignatureKey: config.PayloadSignatureKey,
	}

	// we only need to send the latest metrics,
	// it's no use if the poller puts 10 metrics into the channel, we need only the latest
	// there are 3 possible cases
	// (1) config.ReportInterval == config.PollInterval
	// in this case we can use a buffered channel with cap 1
	//
	// (2) config.ReportInterval > config.PollInterval
	// in this case we can also use a buffered channel with cap 1,
	// but we need to take-and-put new metrics when polling
	//
	// (3) config.ReportInterval < config.PollInterval
	// this means we report the same metrics multiple times
	// it means that we need to take-and-put the same metrics when reporting
	metricsChannel := make(chan metric.MetricsSet, 1)

	go pollingRoutine(ctx, metricsChannel, &pollerInstance, config)
	go reportingRoutine(ctx, metricsChannel, &senderInstance, &metricsStorageClient, config)
}

func pollingRoutine(ctx context.Context, metricsChannel chan metric.MetricsSet, pollerInstance *poller.PollMaker, config *AgentConfig) {
	logger.ZapSugarLogger.Infoln("polling started")

	for i := 0; ; i += 1 {
		select {
		case <-ctx.Done():
			logger.ZapSugarLogger.Infoln("poll context finished")
			return
		default:
			if !pollerInstance.IsRunning {

				if len(metricsChannel) == 0 {
					// if empty, just store in channel
					metricsChannel <- *pollerInstance.Poll()
				} else {
					// take latest and replace it with a new poll
					// use ok not to trigger vet
					_, ok := <-metricsChannel
					if ok {
						metricsChannel <- *pollerInstance.Poll()
					} else {
						metricsChannel <- *pollerInstance.Poll()
					}
				}
			}

			time.Sleep(time.Duration(config.PollInterval) * time.Second)
		}
	}
}

func reportingRoutine(ctx context.Context, metricsChannel chan metric.MetricsSet, senderInstance *sender.MetricsSender, metricsStorageClient *client.MetricsStorageClient, config *AgentConfig) {
	logger.ZapSugarLogger.Infoln("reporting started")

	var metricsSet metric.MetricsSet
	for i := 0; ; i += 1 {
		select {
		case <-ctx.Done():
			logger.ZapSugarLogger.Infoln("report context finished")
			return
		default:

			if !senderInstance.IsRunning {

				if len(metricsChannel) == 0 {
					// nothing to report yet, need to wait for poller
					continue
				} else {
					metricsSet = <-metricsChannel
					metricsChannel <- metricsSet

					senderInstance.Report(&metricsSet, metricsStorageClient)
				}
			}

			time.Sleep(time.Duration(config.ReportInterval) * time.Second)
		}
	}
}
