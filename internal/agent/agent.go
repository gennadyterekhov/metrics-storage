package agent

import (
	"fmt"
	"github.com/gennadyterekhov/metrics-storage/internal/agent/poller"
	"github.com/gennadyterekhov/metrics-storage/internal/agent/sender"
	"runtime"
	"sync"
)

func Agent(address string, reportInterval int, pollInterval int) (err error) {
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

	fmt.Println("pollerInstance.IsRunning", pollerInstance.IsRunning)
	for {
		if !pollerInstance.IsRunning {
			// start periodic poll in bg
			memStatsPtr = pollerInstance.Poll()
			//go pollRoutine(pollInterval, memStatsPtr, metricsChannel)
		}

		if !senderInstance.IsRunning {
			// start periodic send in bg
			senderInstance.Report(memStatsPtr)
			//go reportRoutine(reportInterval, address, metricsChannel)
		}
	}
}
