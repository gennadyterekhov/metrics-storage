package poller

import (
	"fmt"
	"github.com/gennadyterekhov/metrics-storage/internal/agent/metric"
	"github.com/gennadyterekhov/metrics-storage/internal/constants/types"
	"github.com/gennadyterekhov/metrics-storage/internal/logger"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"

	"math/rand"
	"runtime"
)

type PollMaker struct {
	MetricsSet *metric.MetricsSet
	Interval   int
	IsRunning  bool
}

func (pmk *PollMaker) Poll() *metric.MetricsSet {
	pmk.IsRunning = true

	go pmk.saveAdditionalMetrics()

	runtimeStats := &runtime.MemStats{}
	runtime.ReadMemStats(runtimeStats)

	pmk.saveRuntimeStatsToMetricSet(runtimeStats)
	pmk.setNames()
	pmk.setTypes()

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

func (pmk *PollMaker) setNames() {
	pmk.MetricsSet.Alloc.Name = "Alloc"
	pmk.MetricsSet.BuckHashSys.Name = "BuckHashSys"
	pmk.MetricsSet.Frees.Name = "Frees"
	pmk.MetricsSet.GCCPUFraction.Name = "GCCPUFraction"
	pmk.MetricsSet.GCSys.Name = "GCSys"
	pmk.MetricsSet.HeapAlloc.Name = "HeapAlloc"
	pmk.MetricsSet.HeapIdle.Name = "HeapIdle"
	pmk.MetricsSet.HeapInuse.Name = "HeapInuse"
	pmk.MetricsSet.HeapObjects.Name = "HeapObjects"
	pmk.MetricsSet.HeapReleased.Name = "HeapReleased"
	pmk.MetricsSet.HeapSys.Name = "HeapSys"
	pmk.MetricsSet.LastGC.Name = "LastGC"
	pmk.MetricsSet.Lookups.Name = "Lookups"
	pmk.MetricsSet.MCacheInuse.Name = "MCacheInuse"
	pmk.MetricsSet.MCacheSys.Name = "MCacheSys"
	pmk.MetricsSet.MSpanInuse.Name = "MSpanInuse"
	pmk.MetricsSet.MSpanSys.Name = "MSpanSys"
	pmk.MetricsSet.Mallocs.Name = "Mallocs"
	pmk.MetricsSet.NextGC.Name = "NextGC"
	pmk.MetricsSet.NumForcedGC.Name = "NumForcedGC"
	pmk.MetricsSet.NumGC.Name = "NumGC"
	pmk.MetricsSet.OtherSys.Name = "OtherSys"
	pmk.MetricsSet.PauseTotalNs.Name = "PauseTotalNs"
	pmk.MetricsSet.StackInuse.Name = "StackInuse"
	pmk.MetricsSet.StackSys.Name = "StackSys"
	pmk.MetricsSet.Sys.Name = "Sys"
	pmk.MetricsSet.TotalAlloc.Name = "TotalAlloc"
	pmk.MetricsSet.PollCount.Name = "PollCount"
	pmk.MetricsSet.RandomValue.Name = "RandomValue"
}

func (pmk *PollMaker) setTypes() {
	pmk.MetricsSet.Alloc.Type = types.Gauge
	pmk.MetricsSet.BuckHashSys.Type = types.Gauge
	pmk.MetricsSet.Frees.Type = types.Gauge
	pmk.MetricsSet.GCCPUFraction.Type = types.Gauge
	pmk.MetricsSet.GCSys.Type = types.Gauge
	pmk.MetricsSet.HeapAlloc.Type = types.Gauge
	pmk.MetricsSet.HeapIdle.Type = types.Gauge
	pmk.MetricsSet.HeapInuse.Type = types.Gauge
	pmk.MetricsSet.HeapObjects.Type = types.Gauge
	pmk.MetricsSet.HeapReleased.Type = types.Gauge
	pmk.MetricsSet.HeapSys.Type = types.Gauge
	pmk.MetricsSet.LastGC.Type = types.Gauge
	pmk.MetricsSet.Lookups.Type = types.Gauge
	pmk.MetricsSet.MCacheInuse.Type = types.Gauge
	pmk.MetricsSet.MCacheSys.Type = types.Gauge
	pmk.MetricsSet.MSpanInuse.Type = types.Gauge
	pmk.MetricsSet.MSpanSys.Type = types.Gauge
	pmk.MetricsSet.Mallocs.Type = types.Gauge
	pmk.MetricsSet.NextGC.Type = types.Gauge
	pmk.MetricsSet.NumForcedGC.Type = types.Gauge
	pmk.MetricsSet.NumGC.Type = types.Gauge
	pmk.MetricsSet.OtherSys.Type = types.Gauge
	pmk.MetricsSet.PauseTotalNs.Type = types.Gauge
	pmk.MetricsSet.StackInuse.Type = types.Gauge
	pmk.MetricsSet.StackSys.Type = types.Gauge
	pmk.MetricsSet.Sys.Type = types.Gauge
	pmk.MetricsSet.TotalAlloc.Type = types.Gauge
	pmk.MetricsSet.PollCount.Type = types.Counter
	pmk.MetricsSet.RandomValue.Type = types.Gauge
}

func (pmk *PollMaker) saveAdditionalMetrics() {
	memoryStats, err := mem.VirtualMemory()
	if err != nil {
		logger.ZapSugarLogger.Debugln("error when getting psutil stats", err.Error())
		return
	}
	pmk.saveTotalMemory(memoryStats)
	pmk.saveFreeMemory(memoryStats)
	pmk.saveCPUUtilization(memoryStats)
}

func (pmk *PollMaker) saveTotalMemory(memoryStats *mem.VirtualMemoryStat) {
	pmk.MetricsSet.TotalMemory.Value = float64(memoryStats.Total)
	pmk.MetricsSet.TotalMemory.Name = "TotalMemory"
	pmk.MetricsSet.TotalMemory.Type = types.Gauge
}

func (pmk *PollMaker) saveFreeMemory(memoryStats *mem.VirtualMemoryStat) {
	pmk.MetricsSet.FreeMemory.Value = float64(memoryStats.Free)
	pmk.MetricsSet.FreeMemory.Name = "FreeMemory"
	pmk.MetricsSet.FreeMemory.Type = types.Gauge
}

func (pmk *PollMaker) saveCPUUtilization(memoryStats *mem.VirtualMemoryStat) {
	cpus, err := cpu.Percent(0, true)
	pmk.MetricsSet.CPUUtilization = make([]metric.GaugeMetric, len(cpus))
	if err != nil {
		logger.ZapSugarLogger.Debugln("error when getting psutil/cpu stats", err.Error())
		return
	}
	for i := range cpus {
		pmk.MetricsSet.CPUUtilization[i].Value = float64(memoryStats.Used)
		pmk.MetricsSet.CPUUtilization[i].Name = fmt.Sprintf("CPUutilization%v", i+1)
		pmk.MetricsSet.CPUUtilization[i].Type = types.Gauge
	}
}
