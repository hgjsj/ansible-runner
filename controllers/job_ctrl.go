package controllers

import (
	"ansible-runner/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetJob handles the retrieval of a job.
func GetJob(c *gin.Context) {
	id := c.Param("id")
	job, err := services.GetJob(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, job)
}
