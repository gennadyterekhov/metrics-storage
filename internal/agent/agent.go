package agent

import (
	"github.com/gennadyterekhov/metrics-storage/internal/agent/poller"
	"github.com/gennadyterekhov/metrics-storage/internal/agent/sender"
	"runtime"
)

type shouldContinueType func(int) bool

func Agent(address string, shouldContinue shouldContinueType, reportInterval int, pollInterval int) (err error) {
	//metricsChannel := make(chan runtime.MemStats)

	memStatsPtr := &runtime.MemStats{}
	runtime.ReadMemStats(memStatsPtr)

	pollerInstance := poller.PollMaker{
		MemStatsPtr: memStatsPtr,
		Interval:    pollInterval,
	}
	senderInstance := sender.MetricsSender{
		Address:  address,
		Interval: reportInterval,
	}

	_ = pollerInstance.Poll()

	_ = senderInstance.Report(memStatsPtr)

	//for {
	//	// start periodic poll in bg
	//	go pollRoutine(pollInterval, memStatsPtr, metricsChannel)
	//	// start periodic send in bg
	//	go reportRoutine(reportInterval, address, metricsChannel)
	//}

	return nil
}
