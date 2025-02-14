package routers

import (
	"ansible-runner/controllers"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// Server initializes the Gin router and sets up the routes.
func Server() {
	r := gin.Default()
	r.GET("/playbook", controllers.ListPlaybooks)
	r.POST("/playbook/launch", controllers.LaunchPlaybook)
	r.GET("/job/:id", controllers.GetJob)
	port := strconv.Itoa(viper.GetInt("server.port"))
	r.Run(":" + port) // listen and serve on configured port
}
