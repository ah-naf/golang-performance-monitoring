package metrics

import (
	"os"
	"runtime"
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

	// Distribution of Go GC pause durations (in seconds)
	GoGCPauseDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: "myapp",
		Subsystem: "go_runtime",
		Name:      "gc_pause_duration_seconds",
		Help:      "Distribution of Go garbage-collection pause durations in seconds",
		// e.g. 1µs → 10ms → 100ms → 1s → 10s (tweak as you like)
		Buckets: prometheus.ExponentialBuckets(1e-6, 10, 7),
	})
)

func RunSystemMetrics(interval time.Duration) {
	proc, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		panic("metrics: cannot create process handle: " + err.Error())
	}

	if ms, err := proc.CreateTime(); err == nil {
		// 1) Process Start time
		ProcessStartTimeSeconds.Set(float64(ms) / 1000.0)
	}

	var lastNumGC uint32

	go func() {
		for {

			// 2) Open FDs
			if fds, err := proc.NumFDs(); err == nil {
				ProcessOpenFDs.Set(float64(fds))
			}

			// 3) Thread count
			if tn, err := proc.NumThreads(); err == nil {
				ProcessNumThreads.Set(float64(tn))
			}

			// 4) GC pause durations
			var mem runtime.MemStats
			runtime.ReadMemStats(&mem)
			// mem.NumGC is total # of GCs since program start
			// for each new GC, observe its pause
			for i := lastNumGC; i < mem.NumGC; i++ {
				idx := i % uint32(len(mem.PauseNs))
				pauseSec := float64(mem.PauseNs[idx]) / 1e9
				GoGCPauseDuration.Observe(pauseSec)
			}
			lastNumGC = mem.NumGC

			time.Sleep(interval)
		}
	}()
}
