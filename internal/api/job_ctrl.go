package api

import (
	"ansible-runner/internal/services"
	"ansible-runner/internal/tasks"
	"ansible-runner/internal/utils"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// JobResponse represents a job response
type JobResponse struct {
	ID       string `json:"id,omitempty"`
	Playbook string `json:"playbook,omitempty"`
	Status   string `json:"status,omitempty"`
	Step     string `json:"step,omitempty"`
	Output   string `json:"stdout,omitempty"`
	Err      string `json:"stderr,omitempty"`
	StartAt  string `json:"start_at,omitempty"`
	EndAt    string `json:"end_at,omitempty"`
}

// GetJob handles the retrieval of a job.
// @Summary Get Job details
// @Description Get details of a job by ID
// @Tags job
// @Produce json
// @Param id path string true "Job ID"
// @Success 200 {object} JobResponse
// @Failure 500 {object} map[string]string
// @Router /job/{id} [get]
func GetJob(c *gin.Context) {
	id := c.Param("id")
	job, err := services.GetJob(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	utils.Log.Info(fmt.Sprintf("Get job %s and response", job.ID))
	var responseJob = JobResponse{ID: job.ID, Playbook: job.Playbook, Status: job.Status,
		Step: job.Step, Output: job.Output, Err: job.Error, StartAt: job.StartTime.Format(time.RFC3339)}
	if job.EndTime.IsZero() == false {
		responseJob.EndAt = job.EndTime.Format(time.RFC3339)
	}
	c.JSON(http.StatusOK, responseJob)
	gin.Logger()
}

// ListJobs handles the request to list jobs by job name.
// @Summary List all Jobs	
// @Description Get a list of all jobs
// @Tags job
// @Produce json
// @Success 200 {array} JobResponse
// @Failure 500 {object} map[string]string
// @Router /job [get]
func ListJobs(c *gin.Context) {
	j, err := services.ListJobs()

	if err != nil {
		utils.Log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	jobList := make([]JobResponse, 0, len(j))

	for _, job := range j {
		jobList = append(jobList, JobResponse{ID: job.ID, Playbook: job.Playbook, Status: job.Status})
	}
	utils.Log.Info("List inventory is ok and response")
	c.JSON(http.StatusOK, jobList)
}

// DeleteJob handles the request to delete a job.
// @Summary Delete a Job
// @Description Delete a job by ID
// @Tags job
// @Param id path string true "Job ID"
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /job/{id} [delete]
func DeleteJob(c *gin.Context) {
	id := c.Param("id")
	err := services.DeleteJob(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	utils.Log.Info(fmt.Sprintf("Delete job %s and response", id))
	c.JSON(http.StatusOK, gin.H{"status": "Deleted"})
}

// RelaunchJob handles the request to relaunch a job.
// @Summary Relaunch a Job
// @Description Relaunch a job by ID
// @Tags job
// @Param id path string true "Job ID"
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /job/{id}/relaunch [post]
func RelaunchJob(c *gin.Context) {
	id := c.Param("id")

	job, err := services.GetJob(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Check Inventory is exists
	inventory, err := services.GetInventory(job.InventoryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Inventory %s is not exists with error %s", job.ID, err.Error()),
		})
		return
	}

	job.Reset()

	if err := services.SaveJob(job); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.Log.Info(fmt.Sprintf("Relaunch job %s and response", job.ID))
	_, err = tasks.EnqueueJobTask(job, inventory)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"id": job.ID})
}