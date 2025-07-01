package routers

import (
	"database/sql"
	"go-metrics-monitoring-system/internals/database"
	"go-metrics-monitoring-system/internals/models"
	"log"
	"net/http"
	"strconv"

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

func GetAllTask(c *gin.Context) {
	taskRows, err := database.DB.Query("SELECT * FROM Task")
	if err != nil {
		log.Println("Error fetching all tasks from db:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var tasks []models.Task
	for taskRows.Next() {
		var task models.Task
		if err := taskRows.Scan(&task.ID, &task.Title, &task.Description, &task.Status); err != nil {
			log.Println("Error scanning task from rows:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		tasks = append(tasks, task)
	}

	c.JSON(http.StatusOK, gin.H{
		"tasks": tasks,
	})
}

func GetTaskWithID(ctx *gin.Context) {
	// parse the :id param
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
		return
	}

	var task models.Task
	row := database.DB.QueryRow(
		`SELECT id, title, "desc", "status" FROM Task WHERE id = $1`,
		id,
	)
	if err := row.Scan(&task.ID, &task.Title, &task.Description, &task.Status); err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		} else {
			log.Println("Error scanning task:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"task": task})
}

func DeleteTask(ctx *gin.Context) {
    id, err := strconv.Atoi(ctx.Param("id"))
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
        return
    }

    res, err := database.DB.Exec(`DELETE FROM Task WHERE id = $1`, id)
    if err != nil {
        log.Println("Error deleting task:", err)
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    if rows, _ := res.RowsAffected(); rows == 0 {
        ctx.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"message": "task deleted successfully"})
}

func EditTask(ctx *gin.Context) {
    id, err := strconv.Atoi(ctx.Param("id"))
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
        return
    }

    var input models.Task
    if err := ctx.ShouldBind(&input); err != nil {
        log.Printf("Error binding JSON: %v", err)
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    res, err := database.DB.Exec(
        `UPDATE Task
         SET title = $1,
             "desc" = $2,
             "status" = $3
         WHERE id = $4`,
        input.Title,
        input.Description,
        input.Status,
        id,
    )
    if err != nil {
        log.Println("Error updating task:", err)
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    if rows, _ := res.RowsAffected(); rows == 0 {
        ctx.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
        return
    }

    // reflect the correct ID in the response
    input.ID = id
    ctx.JSON(http.StatusOK, gin.H{
        "success": "task updated successfully",
        "task":    input,
    })
}