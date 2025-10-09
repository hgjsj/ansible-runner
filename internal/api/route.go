package api

import (
	"net/http"
	"ansible-runner/internal/docs"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
	swaggerFiles "github.com/swaggo/files" // swagger embed files
)

func SetRoutes() http.Handler {
	r := gin.Default()
	docs.SwaggerInfo.BasePath = "/"
	apiGroup := r.Group("")
	apiGroup.Use(gin.Logger())
	apiGroup.Use(AuthMiddleware)
	apiGroup.GET("/playbook", ListPlaybooks)
	apiGroup.POST("/playbook/launch", LaunchPlaybook)
	apiGroup.GET("/job/:id", GetJob)
	apiGroup.GET("/job", ListJobs)
	apiGroup.DELETE("/job/:id", DeleteJob)
	apiGroup.POST("/job/:id/relaunch", RelaunchJob)
	apiGroup.GET("/inventory", ListInventory)
	apiGroup.GET("/inventory/:id", GetInventory)
	apiGroup.DELETE("/inventory/:id", DeleteInventory)
	apiGroup.POST("/inventory", CreateInventory)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
