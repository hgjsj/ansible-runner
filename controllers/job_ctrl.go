package controllers

import (
	"ansible-runner/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// JobResponse represents a job response
type JobResponse struct {
	ID     string 	`json:"id"`
	Status string	`json:"status"`
	Step   string   `json:"step"`
	Output string   `json:"stdout"`
	Err    string   `json:"stderr"`
}

// GetJob handles the retrieval of a job.
func GetJob(c *gin.Context) {
	id := c.Param("id")
	job, err := services.GetJob(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var responseJob = JobResponse{Status: job.Status, Step: job.Step, Output: job.Output, Err: job.Error}
	c.JSON(http.StatusOK, responseJob)
}
