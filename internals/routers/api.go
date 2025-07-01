package routers

import (
	"go-metrics-monitoring-system/internals/database"
	"go-metrics-monitoring-system/internals/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddNewTask(ctx *gin.Context) {
	var task models.Task
	if err := ctx.ShouldBind(&task); err != nil {
		log.Printf("Error binding JSON: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := database.DB.Exec(`INSERT INTO Task(id, title, "desc", "status") VALUES ($1, $2, $3, $4)`, task.ID, task.Title, task.Description, task.Status)
	if err != nil {
		log.Println("Error inserting new task in db:", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"success": "Task added successfully", "task": task})
}
