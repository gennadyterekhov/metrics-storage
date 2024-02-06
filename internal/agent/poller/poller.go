package poller

import (
	"runtime"
	"time"
)

type Poller interface {
	Poll() error
	wait()
	shouldContinue(int)
}

type PollMaker struct {
	MemStatsPtr *runtime.MemStats
	Interval    int
}

func (pmk *PollMaker) shouldContinue(iter int) bool {
	//return iter == 0
	return true
}

func (pmk *PollMaker) wait() {
	time.Sleep(time.Duration(pmk.Interval * int(time.Second)))
}

func (pmk *PollMaker) pollRoutine(metricsChannel chan runtime.MemStats) {
	pmk.wait()

	runtime.ReadMemStats(pmk.MemStatsPtr)
	//metricsChannel <- *pmk.MemStatsPtr
}

func (pmk *PollMaker) Poll() (err error) {
	metricsChannel := make(chan runtime.MemStats)

	pmk.MemStatsPtr = &runtime.MemStats{}
	runtime.ReadMemStats(pmk.MemStatsPtr)

	for i := 0; pmk.shouldContinue(i); i += 1 {
		pmk.pollRoutine(metricsChannel)
	}

	return nil
}
