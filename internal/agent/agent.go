package agent

import (
	"github.com/gennadyterekhov/metrics-storage/internal/agent/metric"
	"github.com/gennadyterekhov/metrics-storage/internal/agent/poller"
	"github.com/gennadyterekhov/metrics-storage/internal/agent/sender"
	"time"
)

func RunAgent(address string, reportInterval int, pollInterval int) (err error) {
	metricsSet := &metric.MetricsSet{}

	pollerInstance := poller.PollMaker{
		MetricsSet: metricsSet,
		Interval:   pollInterval,
		IsRunning:  false,
	}
	senderInstance := sender.MetricsSender{
		Address:   address,
		Interval:  reportInterval,
		IsRunning: false,
	}

	for i := 0; ; i += 1 {
		time.Sleep(time.Second)

		if !pollerInstance.IsRunning && i%pollInterval == 0 {
			metricsSet = pollerInstance.Poll()
		}

		if !senderInstance.IsRunning && i%reportInterval == 0 {
			senderInstance.Report(metricsSet)
		}
	}
}
