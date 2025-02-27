package controllers

import (
	"ansible-runner/models"
	"ansible-runner/services"
	"ansible-runner/utils"
	"fmt"
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
)

// JobResponse represents a job response
type JobResponse struct {
	ID     string 	`json:"id,omitempty"`
	Status string	`json:"status,omitempty"`
	Step   string   `json:"step,omitempty"`
	Output string   `json:"stdout,omitempty"`
	Err    string   `json:"stderr,omitempty"`
	Inventory *models.Inventory `json:"inventory,omitempty"`
	StartAt string `json:"start_at,omitempty"`
	EndAt string `json:"end_at,omitempty"`
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
	var responseJob = JobResponse{ID: job.ID, Status: job.Status, Step: job.Step, Output: job.Output, Err: job.Error,
		StartAt: job.StartTime.Format(time.RFC3339)}
	if job.EndTime.IsZero() == false {
		responseJob.EndAt = job.EndTime.Format(time.RFC3339)
	}
	c.JSON(http.StatusOK, responseJob)
}

func GetJobInventory(c *gin.Context){
	id := c.Param("id")
	inv, err := services.GetJobInventory(id)
	if err != nil {
		utils.Log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	utils.Log.Info(fmt.Sprintf("Get job inventory %s and response", id))
    inv.Localhost = nil
	var responseJob = JobResponse{Inventory: inv}
	c.JSON(http.StatusOK, responseJob)
}
