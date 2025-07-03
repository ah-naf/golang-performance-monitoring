package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func RunMetrics(interval time.Duration) {
	RunCPUMetrics(interval)
	RunMemoryMetrics(interval)
	RunSystemMetrics(interval)
}

func RegisterAllMetrics(reg *prometheus.Registry) {
	reg.MustRegister(
		// CPU metrics
		ProcessCPUSecondsTotal,
		ProcessCPUUsagePercent,

		// System metrics
		ProcessStartTimeSeconds,
		ProcessOpenFDs,
		ProcessNumThreads,
		GoGCPauseDuration,

		// Memory metrics
		ProcessResidentMemoryBytes,
		ProcessVirtualMemoryBytes,

		// HTTP metrics
		HTTPRequestsTotal,
		HTTPRequestDuration,
		HTTPRequestSizeBytes,
		HTTPResponseSizeBytes,
	)
}
