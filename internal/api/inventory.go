package api

import (
	"ansible-runner/internal/db"
	"ansible-runner/internal/models"
	"ansible-runner/internal/services"
	"ansible-runner/internal/utils"
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
	requestBody := models.NewInventory(false)
	if err := c.ShouldBindJSON(requestBody); err != nil {
		utils.Log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}
	if _, err := services.CreateInventory(requestBody); err != nil {
		utils.Log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid request body",
		})
		return
	}
	utils.Log.Info("Create inventory is ok and response")
	c.JSON(http.StatusCreated, requestBody)
}

func DeleteInventory(c *gin.Context) {
	id := c.Param("id")
	err := services.DeleteInventory(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	utils.Log.Info(fmt.Sprintf("Delete job %s and response", id))
	c.JSON(http.StatusOK, gin.H{"status": "Deleted"})
}