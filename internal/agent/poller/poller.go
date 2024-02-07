package poller

import (
	"github.com/gennadyterekhov/metrics-storage/internal/agent/metric"
	"math/rand"
	"runtime"
	"time"
)

type PollMaker struct {
	MetricsSet *metric.MetricsSet
	Interval   int
	IsRunning  bool
}

func (pmk *PollMaker) wait() {
	time.Sleep(time.Duration(pmk.Interval * int(time.Second)))
}

func (pmk *PollMaker) Poll() *metric.MetricsSet {
	pmk.IsRunning = true
	pmk.wait()

	runtimeStats := &runtime.MemStats{}
	runtime.ReadMemStats(runtimeStats)

	pmk.saveRuntimeStatsToMetricSet(runtimeStats)
	pmk.IsRunning = false

	return pmk.MetricsSet
}

func (pmk *PollMaker) saveRuntimeStatsToMetricSet(runtimeStats *runtime.MemStats) {
	pmk.MetricsSet.Alloc.Value = float64(runtimeStats.Alloc)
	pmk.MetricsSet.BuckHashSys.Value = float64(runtimeStats.BuckHashSys)
	pmk.MetricsSet.Frees.Value = float64(runtimeStats.Frees)
	pmk.MetricsSet.GCCPUFraction.Value = float64(runtimeStats.GCCPUFraction)
	pmk.MetricsSet.GCSys.Value = float64(runtimeStats.GCSys)
	pmk.MetricsSet.HeapAlloc.Value = float64(runtimeStats.HeapAlloc)
	pmk.MetricsSet.HeapIdle.Value = float64(runtimeStats.HeapIdle)
	pmk.MetricsSet.HeapInuse.Value = float64(runtimeStats.HeapInuse)
	pmk.MetricsSet.HeapObjects.Value = float64(runtimeStats.HeapObjects)
	pmk.MetricsSet.HeapReleased.Value = float64(runtimeStats.HeapReleased)
	pmk.MetricsSet.HeapSys.Value = float64(runtimeStats.HeapSys)
	pmk.MetricsSet.LastGC.Value = float64(runtimeStats.LastGC)
	pmk.MetricsSet.Lookups.Value = float64(runtimeStats.Lookups)
	pmk.MetricsSet.MCacheInuse.Value = float64(runtimeStats.MCacheInuse)
	pmk.MetricsSet.MCacheSys.Value = float64(runtimeStats.MCacheSys)
	pmk.MetricsSet.MSpanInuse.Value = float64(runtimeStats.MSpanInuse)
	pmk.MetricsSet.MSpanSys.Value = float64(runtimeStats.MSpanSys)
	pmk.MetricsSet.Mallocs.Value = float64(runtimeStats.Mallocs)
	pmk.MetricsSet.NextGC.Value = float64(runtimeStats.NextGC)
	pmk.MetricsSet.NumForcedGC.Value = float64(runtimeStats.NumForcedGC)
	pmk.MetricsSet.NumGC.Value = float64(runtimeStats.NumGC)
	pmk.MetricsSet.OtherSys.Value = float64(runtimeStats.OtherSys)
	pmk.MetricsSet.PauseTotalNs.Value = float64(runtimeStats.PauseTotalNs)
	pmk.MetricsSet.StackInuse.Value = float64(runtimeStats.StackInuse)
	pmk.MetricsSet.StackSys.Value = float64(runtimeStats.StackSys)
	pmk.MetricsSet.Sys.Value = float64(runtimeStats.Sys)
	pmk.MetricsSet.TotalAlloc.Value = float64(runtimeStats.TotalAlloc)
	pmk.MetricsSet.PollCount.Value += 1
	pmk.MetricsSet.RandomValue.Value = rand.Float64()
}
