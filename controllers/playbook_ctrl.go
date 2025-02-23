package controllers

import (
	"ansible-runner/models"
	"ansible-runner/services"
	"ansible-runner/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// LaunchRequest represents a request to launch a playbook.
type LaunchRequest struct {
	Playbook  string           `json:"playbook"`
	Inventory models.Inventory `json:"inventory"`
}

// ListPlaybooks handles the listing of playbooks.
func ListPlaybooks(c *gin.Context) {
	baseDir := viper.GetString("playbook.baseDir")
	playbooks, err := services.GetPlaybooks(baseDir)
	if err != nil {
		utils.Log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	utils.Log.Info("List playbook is ok and response")
	c.JSON(http.StatusOK, gin.H{
		"playbooks": playbooks,
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

	utils.Log.Info("Get payload from request is ok")
	job, err := services.RunPlaybook(requestBody.Playbook, requestBody.Inventory)
	if err != nil {
		utils.Log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	utils.Log.Info("Run playbook is ok and response")
	var responseJob = JobResponse{ID: job.ID,Status: job.Status}
	c.JSON(http.StatusOK,responseJob)
}
