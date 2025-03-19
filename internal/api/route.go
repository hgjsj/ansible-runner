package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetRoutes() http.Handler {
	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(AuthMiddleware)
	r.GET("/playbook", ListPlaybooks)
	r.POST("/playbook/launch", LaunchPlaybook)
	r.GET("/job/:id", GetJob)
	r.GET("/job", ListJobs)
	r.GET("/inventory", ListInventory)
	r.GET("/inventory/:id", GetInventory)
	r.POST("/inventory", CreateInventory)
	return r
}
