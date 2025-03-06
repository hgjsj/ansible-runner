package routers

import (
	"ansible-runner/controllers"
	"ansible-runner/utils"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var srv *http.Server


// Server initializes the Gin router and sets up the routes.
func Server() {
	r := gin.Default()
	r.Use(gin.Logger())
	r.GET("/playbook", controllers.ListPlaybooks)
	r.POST("/playbook/launch", controllers.LaunchPlaybook)
	r.GET("/job/:id", controllers.GetJob)
	r.GET("/job", controllers.ListJobs)
	r.GET("/inventory", controllers.ListInventory)
	r.GET("/inventory/:id", controllers.GetInventory)
	r.POST("/inventory", controllers.CreateInventory)
	port := strconv.Itoa(viper.GetInt("server.port"))
	utils.Log.Info(fmt.Sprintf("Server start and listen on port %s", port))
	srv = new(http.Server)
	srv.Addr =  fmt.Sprintf(":%s", port)
	srv.Handler = r

	go func(){
		if err := srv.ListenAndServeTLS("./cert/server.crt", "./cert/server.key"); err != nil && err != http.ErrServerClosed {
			fmt.Printf("listen: %s\n", err)
		}
	}()
}

func ShutDown(){
	if srv != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		srv.Shutdown(ctx)
	}
}
