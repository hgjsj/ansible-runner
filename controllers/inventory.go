package controllers

import (
	"ansible-runner/db"
	"ansible-runner/models"
	"ansible-runner/services"
	"ansible-runner/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ListInventory(c *gin.Context) {
	i, e := services.ListInventory(db.NewMongoClient().InventoryCollection)
	if e != nil {
		utils.Log.Error(e.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": e.Error()})
		return
	}
	utils.Log.Info("List inventory is ok and response")
	c.JSON(http.StatusOK, i)
}

func GetInventory(c *gin.Context) {
	id := c.Param("id")
	inv, err := services.GetInventory(id)
	if err != nil {
		utils.Log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	utils.Log.Info(fmt.Sprintf("Get inventory %s and response", id))
	c.JSON(http.StatusOK, inv)
}

func CreateInventory(c *gin.Context) {
	var requestBody models.Inventory
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		utils.Log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}
	if err := services.CreateInventory(&requestBody); err != nil {
		utils.Log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid request body",
		})
		return
	}
	utils.Log.Info("Create inventory is ok and response")
	c.JSON(http.StatusOK, requestBody)
}
