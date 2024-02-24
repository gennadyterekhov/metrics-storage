package agent

import (
	"github.com/gennadyterekhov/metrics-storage/internal/agent/healthcheck"
	"github.com/gennadyterekhov/metrics-storage/internal/agent/metric"
	"github.com/gennadyterekhov/metrics-storage/internal/agent/poller"
	"github.com/gennadyterekhov/metrics-storage/internal/agent/sender"
	"github.com/gennadyterekhov/metrics-storage/internal/logger"
	"time"
)

type AgentConfig struct {
	Addr           string
	ReportInterval int
	PollInterval   int
}

func RunAgent(config *AgentConfig) (err error) {
	metricsSet := &metric.MetricsSet{}

	pollerInstance := poller.PollMaker{
		MetricsSet: metricsSet,
		Interval:   config.PollInterval,
		IsRunning:  false,
	}
	senderInstance := sender.MetricsSender{
		Address:   config.Addr,
		Interval:  config.ReportInterval,
		IsRunning: false,
	}

	for i := 0; ; i += 1 {

		// TODO fix interval issue
		// У нас pollInterval 2с, reportInterval 10с
		// Какой будет метрика PollCount на сервере через 20с?
		// Условно мы 10 раз сделали poll и 2 раза репорт.
		// В идеальном мире(все операции моментальны) она должна бы быть равна 10, а будет?
		time.Sleep(time.Second)

		if !pollerInstance.IsRunning && i%config.PollInterval == 0 {
			metricsSet = pollerInstance.Poll()
		}

		if !senderInstance.IsRunning && i%config.ReportInterval == 0 {
			if !isServerAvailable(config) {
				logger.ZapSugarLogger.Warnln(
					"agent will not send metrics because server healthcheck was not successful",
				)
				continue
			}
			senderInstance.Report(metricsSet)
		}
	}
}

func isServerAvailable(config *AgentConfig) (isAvailable bool) {
	return healthcheck.MakeHealthcheck(config.Addr)
}
