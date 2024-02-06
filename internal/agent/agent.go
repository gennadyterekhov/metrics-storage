package agent

import (
	"github.com/gennadyterekhov/metrics-storage/internal/agent/poller"
	"github.com/gennadyterekhov/metrics-storage/internal/agent/sender"
	"runtime"
	"sync"
)

type shouldContinueType func(int) bool

func Agent(address string, shouldContinue shouldContinueType, reportInterval int, pollInterval int) (err error) {
	metricsChannel := make(chan runtime.MemStats)

	memStatsPtr := &runtime.MemStats{}
	runtime.ReadMemStats(memStatsPtr)

	isRoutineRunningMutex := &sync.Mutex{}
	pollerInstance := poller.PollMaker{
		MemStatsPtr: memStatsPtr,
		Interval:    pollInterval,
		Channel:     metricsChannel,
		IsRunning:   false,
		IsRunningMu: isRoutineRunningMutex,
	}
	senderInstance := sender.MetricsSender{
		Address:     address,
		Interval:    reportInterval,
		Channel:     metricsChannel,
		IsRunning:   false,
		IsRunningMu: isRoutineRunningMutex,
	}

	for {
		if !pollerInstance.IsRunning {
			// start periodic poll in bg
			go pollerInstance.Poll()
			//go pollRoutine(pollInterval, memStatsPtr, metricsChannel)
		}

		if !senderInstance.IsRunning {
			// start periodic send in bg
			go senderInstance.Report()
			//go reportRoutine(reportInterval, address, metricsChannel)
		}
	}
}
