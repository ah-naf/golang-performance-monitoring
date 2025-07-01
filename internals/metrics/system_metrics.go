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
)

func RunCPUMetrics(interval time.Duration) {
	proc, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		panic("metrics: cannot create process handle: " + err.Error())
	}

	var lastTotal float64

	go func() {
		for {
			// 3a) CPU seconds: user + system
			if times, err := proc.Times(); err == nil {
				total := times.User + times.System
				delta := total - lastTotal
				if delta < 0 {
					delta = 0
				}
				ProcessCPUSecondsTotal.Add(delta)
				lastTotal = total
			}

			// 3b) CPU usage % over the last interval
			if pct, err := proc.Percent(0); err == nil {
				ProcessCPUUsagePercent.Set(pct)
			}

			time.Sleep(interval)
		}
	}()
}
