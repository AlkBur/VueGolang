package handlers

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"net/http"
	"github.com/AlkBur/VueGolang/models"
	"database/sql"
)

// GetTasks endpoint
func GetTasks(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, models.GetTasks(db))
	}
}

// PutTask endpoint
func PutTask(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Instantiate a new task
		var task models.Task
		// Map imcoming JSON body to the new Task
		c.Bind(&task)
		// Add a task using our new model
		id, err := models.PutTask(db, task.Name)
		// Return a JSON response if successful
		if err == nil {
			c.JSON(http.StatusCreated, gin.H{
				"created": id,
			})
			// Handle any errors
		}else{
			c.JSON(422, gin.H{"error": err.Error()})
		}
	}
}
// DeleteTask endpoint
func DeleteTask(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		// Use our new model to delete a task
		_, err := models.DeleteTask(db, id)
		// Return a JSON response on success
		if err == nil {
			c.JSON(http.StatusOK, gin.H{
				"deleted": id,
			})
			// Handle errors
		}else {
			c.JSON(422, gin.H{"error": err.Error()})
		}
	}
}