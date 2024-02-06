package poller

import (
	"runtime"
	"sync"
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
	IsRunning   bool
	IsRunningMu *sync.Mutex
}

func (pmk *PollMaker) shouldContinue(iter int) bool {
	//return iter == 0
	return true
}

func (pmk *PollMaker) wait() {
	time.Sleep(time.Duration(pmk.Interval * int(time.Second)))
}

func (pmk *PollMaker) pollRoutine() {
	pmk.IsRunningMu.Lock()
	pmk.IsRunning = true
	pmk.wait()

	runtime.ReadMemStats(pmk.MemStatsPtr)
	//fmt.Println("updated runtime metrics, saving to channel")
	pmk.Channel <- *pmk.MemStatsPtr
	pmk.IsRunningMu.Unlock()

}

func (pmk *PollMaker) Poll() {
	pmk.MemStatsPtr = &runtime.MemStats{}
	runtime.ReadMemStats(pmk.MemStatsPtr)

	pmk.pollRoutine()
}
