package controllers

import (
	"ansible-runner/models"
	"ansible-runner/services"
	"ansible-runner/tasks"
	"ansible-runner/utils"
	"fmt"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
)

// LaunchRequest represents a request to launch a playbook.
type LaunchRequest struct {
	Playbook  string `json:"playbook"`
	Inventory string `json:"inventory"`
}

var playbookList []string = []string{"p4s_account_validate.yml"}

// ListPlaybooks handles the listing of playbooks.
func ListPlaybooks(c *gin.Context) {
	utils.Log.Info("List playbook is ok and response")
	c.JSON(http.StatusOK, gin.H{
		"playbooks": playbookList,
	})
}

// LaunchPlaybook handles the launch of a playbook.
func LaunchPlaybook(c *gin.Context) {
	var requestBody LaunchRequest
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		utils.Log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// Check if the playbook is in the list of available playbooks.
	if slices.Contains(playbookList, requestBody.Playbook) == false {
		utils.Log.Error(fmt.Sprintf("Playbook %s is not in the list of available playbooks", requestBody.Playbook))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Playbook %s is not in the list of available playbooks", requestBody.Playbook),
		})
		return
	}

	// save ansible job to db
	aj := models.NewAnsibleJob(true)
	aj.Playbook = requestBody.Playbook
	aj.Status = models.Scheduled
	if err := services.SaveJob(aj); err != nil {
		utils.Log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Send message to queue
	utils.Log.Info("Get payload from request is ok")
	_, err := tasks.EnqueueJobTask(aj.ID.Hex(), requestBody.Inventory)
	if err != nil {
		utils.Log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	utils.Log.Info("Run playbook is ok and response")
	var responseJob = JobResponse{ID: aj.ID.Hex(), Status: aj.Status}
	c.JSON(http.StatusOK, responseJob)
}
