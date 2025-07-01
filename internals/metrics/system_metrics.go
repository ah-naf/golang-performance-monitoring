package metrics

import (
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/v3/process"
)

var (
	// Total CPU time (in seconds) consumed by this process
	ProcessCPUSecondsTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "myapp",
		Subsystem: "process",
		Name:      "cpu_seconds_total",
		Help:      "Total CPU time consumed by the process in seconds",
	})

	// Instantaneous CPU usage percentage of this process
	ProcessCPUUsagePercent = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "myapp",
		Subsystem: "process",
		Name:      "cpu_usage_percent",
		Help:      "Current CPU utilization percentage of the process",
	})

	// UNIX timestamp of when the process started
	ProcessStartTimeSeconds = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "myapp",
		Subsystem: "process",
		Name:      "start_time_seconds",
		Help:      "UNIX timestamp of when the process started",
	})

	// Number of open file descriptors
	ProcessOpenFDs = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "myapp",
		Subsystem: "process",
		Name:      "open_fds",
		Help:      "Number of open file descriptors by the process",
	})

	// Number of OS threads in this process
	ProcessNumThreads = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "myapp",
		Subsystem: "process",
		Name:      "threads",
		Help:      "Number of OS threads in the process",
	})

	// Resident (physical) memory in bytes currently used by the process
	ProcessResidentMemoryBytes = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "myapp",
		Subsystem: "process",
		Name:      "resident_memory_bytes",
		Help:      "Resident (physical) memory in bytes currently used by the process",
	})

	// Virtual memory in bytes allocated by the process (RSS + swap + mappings)
	ProcessVirtualMemoryBytes = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "myapp",
		Subsystem: "process",
		Name:      "virtual_memory_bytes",
		Help:      "Virtual memory in bytes allocated by the process (RSS + swap + mappings)",
	})
)

func RunCPUMetrics(interval time.Duration) {
	proc, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		panic("metrics: cannot create process handle: " + err.Error())
	}

	if ms, err := proc.CreateTime(); err == nil {
		ProcessStartTimeSeconds.Set(float64(ms) / 1000.0)
	}

	var lastTotal float64

	go func() {
		for {
			// 1) CPU seconds: user + system
			if times, err := proc.Times(); err == nil {
				total := times.User + times.System
				delta := total - lastTotal
				if delta < 0 {
					delta = 0
				}
				ProcessCPUSecondsTotal.Add(delta)
				lastTotal = total
			}

			// 2) CPU usage % over the last interval
			if pct, err := proc.Percent(0); err == nil {
				ProcessCPUUsagePercent.Set(pct)
			}

			// 3) Open FDs
			if fds, err := proc.NumFDs(); err == nil {
				ProcessOpenFDs.Set(float64(fds))
			}

			// 4) Thread count
			if tn, err := proc.NumThreads(); err == nil {
				ProcessNumThreads.Set(float64(tn))
			}

			// Resident memory
			if mi, err := proc.MemoryInfo(); err == nil {
				ProcessResidentMemoryBytes.Set(float64(mi.RSS))
			}

			// Virtual memory
			if mi, err := proc.MemoryInfo(); err == nil {
				ProcessVirtualMemoryBytes.Set(float64(mi.VMS))
			}

			time.Sleep(interval)
		}
	}()
}
