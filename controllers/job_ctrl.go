package controllers

import (
	"ansible-runner/services"
	"ansible-runner/utils"
	"fmt"
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
		utils.Log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	utils.Log.Info(fmt.Sprintf("Get job %s and response", job.ID))
	var responseJob = JobResponse{ID: job.ID, Status: job.Status, Step: job.Step, Output: job.Output, Err: job.Error}
	c.JSON(http.StatusOK, responseJob)
}
