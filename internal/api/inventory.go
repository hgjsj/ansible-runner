package api

import (
	"ansible-runner/internal/models"
	"ansible-runner/internal/services"
	"ansible-runner/internal/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ListInventory handles the listing of inventories.
// @Summary List Inventories
// @Description Get a list of available inventories
// @Tags inventory
// @Produce json
// @Success 200 {array} models.Inventory
// @Failure 500 {object} map[string]string
// @Router /inventory [get]
func ListInventory(c *gin.Context) {
	i, e := services.ListInventory()
	if e != nil {
		utils.Log.Error(e.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": e.Error()})
		return
	}
	utils.Log.Info("List inventory is ok and response")
	c.JSON(http.StatusOK, i)
}

// GetInventory handles the retrieval of an inventory.
// @Summary Get Inventory details
// @Description Get details of an inventory by ID
// @Tags inventory
// @Produce json
// @Param id path string true "Inventory ID"
// @Success 200 {object} models.Inventory
// @Failure 500 {object} map[string]string
// @Router /inventory/{id} [get]
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


// CreateInventory handles the creation of a new inventory.
// @Summary Create Inventory
// @Description Create a new inventory
// @Tags inventory
// @Accept json
// @Produce json
// @Param inventory body models.Inventory true "Inventory to create"
// @Success 201 {object} models.Inventory
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /inventory [post]
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

// DeleteInventory handles the deletion of an inventory.
// @Summary Delete Inventory
// @Description Delete an inventory by ID
// @Tags inventory
// @Param id path string true "Inventory ID"
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /inventory/{id} [delete]	
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