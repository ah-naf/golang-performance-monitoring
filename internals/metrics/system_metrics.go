package metrics

import (
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/v3/process"
)

var (
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
)

func RunSystemMetrics(interval time.Duration) {
	proc, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		panic("metrics: cannot create process handle: " + err.Error())
	}

	if ms, err := proc.CreateTime(); err == nil {
		ProcessStartTimeSeconds.Set(float64(ms) / 1000.0)
	}

	go func() {
		for {

			// 3) Open FDs
			if fds, err := proc.NumFDs(); err == nil {
				ProcessOpenFDs.Set(float64(fds))
			}

			// 4) Thread count
			if tn, err := proc.NumThreads(); err == nil {
				ProcessNumThreads.Set(float64(tn))
			}

			time.Sleep(interval)
		}
	}()
}
