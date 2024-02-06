package poller

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

type Poller interface {
	Poll() error
	wait()
}

type PollMaker struct {
	MemStatsPtr *runtime.MemStats
	Interval    int
	Channel     chan runtime.MemStats
	IsRunning   bool
	IsRunningMu *sync.Mutex
}

func (pmk *PollMaker) wait() {
	time.Sleep(time.Duration(pmk.Interval * int(time.Second)))
}

func (pmk *PollMaker) Poll() *runtime.MemStats {
	pmk.IsRunningMu.Lock()
	pmk.IsRunning = true
	fmt.Println("pmk.IsRunning", pmk.IsRunning)
	pmk.wait()

	pmk.MemStatsPtr = &runtime.MemStats{}
	runtime.ReadMemStats(pmk.MemStatsPtr)

	fmt.Println("updated runtime metrics, saving to channel")
	//pmk.Channel <- *pmk.MemStatsPtr
	fmt.Println("SAVED to channel")

	pmk.IsRunning = false

	pmk.IsRunningMu.Unlock()
	return pmk.MemStatsPtr
}
