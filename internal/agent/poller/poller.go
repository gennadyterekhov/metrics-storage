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
	Channel     chan runtime.MemStats
}

func (pmk *PollMaker) shouldContinue(iter int) bool {
	//return iter == 0
	return true
}

func (pmk *PollMaker) wait() {
	time.Sleep(time.Duration(pmk.Interval * int(time.Second)))
}

func (pmk *PollMaker) pollRoutine() {
	pmk.wait()

	runtime.ReadMemStats(pmk.MemStatsPtr)
	pmk.Channel <- *pmk.MemStatsPtr
}

func (pmk *PollMaker) Poll() {
	pmk.MemStatsPtr = &runtime.MemStats{}
	runtime.ReadMemStats(pmk.MemStatsPtr)

	pmk.pollRoutine()
}
