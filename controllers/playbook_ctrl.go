package controllers

import (
	"ansible-runner/models"
	"ansible-runner/services"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"playbooks": playbooks,
	})
}

// LaunchPlaybook handles the launch of a playbook.
func LaunchPlaybook(c *gin.Context) {
	var requestBody LaunchRequest
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	job, err := services.RunPlaybook(requestBody.Playbook, requestBody.Inventory)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"job": job,
	})
}
