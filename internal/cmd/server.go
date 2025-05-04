package cmd

import (
	"ansible-runner/internal/api"
	"ansible-runner/internal/db"
	"ansible-runner/internal/tasks"
	"ansible-runner/internal/utils"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

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
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		srv.Shutdown(ctx)
	}
	tasks.Close()
	tasks.CloseClient()
	db.MongoClient.Close()
	utils.CleanUp()
}

// Server initializes the Gin router and sets up the routes.
func server() {

	port := strconv.Itoa(viper.GetInt("server.port"))
	utils.Log.Info(fmt.Sprintf("Server start and listen on port %s", port))
	srv = new(http.Server)
	srv.Addr = fmt.Sprintf(":%s", port)
	srv.Handler = api.SetRoutes()

	go func() {
		if err := srv.ListenAndServeTLS("./cert/server.crt", "./cert/server.key"); err != nil && err != http.ErrServerClosed {
			utils.Log.Error(err.Error())
			panic(err)
		}
	}()

	go func ()  {
		if err := tasks.LaunchServer("server"); err != nil {
			utils.Log.Error(err.Error())
			panic(err)
		}
	}()
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Launch an HTTP server to accept requests to schedule  Ansible jobs",
	Run: func(cmd *cobra.Command, args []string) {
		utils.LoadConfig("./configs/config")
		utils.InitLogger()
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
