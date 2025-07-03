package metrics

import (
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/v3/process"
)

var (
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

func RunMemoryMetrics(interval time.Duration) {
	proc, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		panic("metrics: cannot create process handle: " + err.Error())
	}

	go func() {
		for {
			// // Resident memory
			if mi, err := proc.MemoryInfo(); err == nil {
				ProcessResidentMemoryBytes.Set(float64(mi.RSS))
			}

			// // Virtual memory
			if mi, err := proc.MemoryInfo(); err == nil {
				ProcessVirtualMemoryBytes.Set(float64(mi.VMS))
			}

			time.Sleep(interval)
		}
	}()
}
