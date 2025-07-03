package metrics

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// Total number of HTTP requests received, partitioned by method, endpoint and status code.
	HTTPRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "myapp",
			Subsystem: "http",
			Name:      "requests_total",
			Help:      "Total number of HTTP requests received.",
		},
		[]string{"method", "endpoint", "status_code"},
	)

	// Histogram of HTTP request durations in seconds, partitioned by method, endpoint and status code.
	HTTPRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "myapp",
			Subsystem: "http",
			Name:      "request_duration_seconds",
			Help:      "Buckets of HTTP request durations (in seconds).",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"method", "endpoint", "status_code"},
	)

	// Histogram of HTTP request payload sizes in bytes, partitioned by method and endpoint.
	HTTPRequestSizeBytes = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "myapp",
			Subsystem: "http",
			Name:      "request_size_bytes",
			Help:      "Buckets of HTTP request payload sizes (in bytes).",
			// e.g. 100B, 1KB, 10KB, 100KB, 1MB, 10MB
			Buckets: prometheus.ExponentialBuckets(100, 10, 6),
		},
		[]string{"method", "endpoint"},
	)

	// Histogram of HTTP response payload sizes in bytes, partitioned by method, endpoint and status code.
	HTTPResponseSizeBytes = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "myapp",
			Subsystem: "http",
			Name:      "response_size_bytes",
			Help:      "Buckets of HTTP response payload sizes (in bytes).",
			Buckets:   prometheus.ExponentialBuckets(100, 10, 6),
		},
		[]string{"method", "endpoint", "status_code"},
	)
)

type bodyWriter struct {
	gin.ResponseWriter
	size int
}

func (w *bodyWriter) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	w.size += n
	return n, err
}

func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		bw := &bodyWriter{ResponseWriter: c.Writer}
		c.Writer = bw

		c.Next()

		duration := time.Since(start).Seconds()
		method := c.Request.Method

		endpoint := c.FullPath()
		if endpoint == "" {
			endpoint = c.Request.URL.Path
		}

		statusCode := strconv.Itoa(c.Writer.Status())

		reqSize := c.Request.ContentLength
		if reqSize < 0 {
			reqSize = 0
		}

		// Observe metrics
		HTTPRequestsTotal.
			WithLabelValues(method, endpoint, statusCode).
			Inc()

		HTTPRequestDuration.
			WithLabelValues(method, endpoint, statusCode).
			Observe(duration)

		HTTPRequestSizeBytes.
			WithLabelValues(method, endpoint).
			Observe(float64(reqSize))

		HTTPResponseSizeBytes.
			WithLabelValues(method, endpoint, statusCode).
			Observe(float64(bw.size))
	}
}
