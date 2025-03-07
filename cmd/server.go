package cmd

import (
	"ansible-runner/configs"
	"ansible-runner/controllers"
	"ansible-runner/db"
	"ansible-runner/utils"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var srv *http.Server

func cleanup(cmd *cobra.Command, args []string) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	utils.Log.Info("Cleanup and shutting down server...")
	if srv != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
		defer cancel()
		srv.Shutdown(ctx)
	}
	db.MongoClient.Close()
	utils.CleanUp()
}


// Server initializes the Gin router and sets up the routes.
func server() {
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

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Launch an HTTP server to accept requests to schedule  Ansible jobs",
	Run: func(cmd *cobra.Command, args []string) {
		utils.InitLogger("./logs/server.log")
		configs.LoadConfig("./")
		server()
	},
	PostRun: cleanup,
}

func init() {
	rootCmd.AddCommand(serverCmd)
	// add a flag to set the port
	serverCmd.Flags().StringP("port", "p", "8080", "Port to listen on")

	// bind the flag to the viper config
	viper.BindPFlag("server.port", serverCmd.Flags().Lookup("port"))
}
