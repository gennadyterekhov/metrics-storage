package agent

import (
	"github.com/gennadyterekhov/metrics-storage/internal/agent/metric"
	"github.com/gennadyterekhov/metrics-storage/internal/agent/poller"
	"github.com/gennadyterekhov/metrics-storage/internal/agent/sender"
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
		time.Sleep(time.Second)

		if !pollerInstance.IsRunning && i%config.PollInterval == 0 {
			metricsSet = pollerInstance.Poll()
		}

		if !senderInstance.IsRunning && i%config.ReportInterval == 0 {
			senderInstance.Report(metricsSet)
		}
	}
}
