package routers

import (
	"fmt"
	"net/http"
	"syscall"

	"go-metrics-monitoring-system/internals/database"

	"github.com/gin-gonic/gin"
)

func HealthCheck(c *gin.Context) {
	// 1) DB ping
	if err := database.DB.Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "error",
			"healthy": false,
			"error":   "db: " + err.Error(),
		})
		return
	}

	// 2) Disk space
	var fs syscall.Statfs_t
	if err := syscall.Statfs("/", &fs); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "error",
			"healthy": false,
			"error":   "disk: " + err.Error(),
		})
		return
	}
	// available blocks * size per block = available bytes
	total := fs.Blocks * uint64(fs.Bsize)
	free := fs.Bavail * uint64(fs.Bsize)
	used := total - free
	usage := float64(used) / float64(total) * 100

	// 3) decide overall health (e.g. fail if > 90% used)
	healthyDisk := usage < 90.0

	// 4) return combined status
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"healthy": healthyDisk,
		"db":      "up",
		"disk": gin.H{
			"total_bytes": total,
			"free_bytes":  free,
			"used_bytes":  used,
			"used_pct":    fmt.Sprintf("%.1f%%", usage),
		},
	})
}
