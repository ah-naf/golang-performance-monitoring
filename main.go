package main

import (
	"context"
	"go-metrics-monitoring-system/internals/database"
	"go-metrics-monitoring-system/internals/metrics"
	"go-metrics-monitoring-system/internals/routers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	router := gin.Default()
	database.ConnectDatabase()

	metrics.RunCPUMetrics(10 * time.Second)

	router.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Golang Metric Monitoring System")
	})

	task := router.Group("/task")
	{
		task.POST("", routers.AddNewTask)
		task.GET("", routers.GetAllTask)
		task.GET("/:id", routers.GetTaskWithID)
		task.PUT("/:id", routers.EditTask)
		task.DELETE("/:id", routers.DeleteTask)
	}
	router.GET("/health", routers.HealthCheck)

	reg := prometheus.NewRegistry()
	reg.MustRegister(
		metrics.ProcessCPUSecondsTotal, metrics.ProcessCPUUsagePercent,
	)

	router.GET("/metrics", gin.WrapH(
		promhttp.HandlerFor(reg, promhttp.HandlerOpts{}),
	))

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router.Handler(),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Println("Server Shutdown:", err)
	}

	<-ctx.Done()
	log.Println("timeout of 1 seconds.")
	log.Println("Server exiting")
}
