package cmd

import (
	"ansible-runner/configs"
	"ansible-runner/db"
	"ansible-runner/routers"
	"ansible-runner/utils"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func cleanup(cmd *cobra.Command, args []string) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	utils.Log.Info("Cleanup and shutting down server...")
	routers.ShutDown()
	db.MongoClient.Close()
	utils.CleanUp()
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Launch an HTTP server to accept requests to schedule  Ansible jobs",
	Run: func(cmd *cobra.Command, args []string) {
		utils.InitLogger()
		configs.LoadConfig("./")
		routers.Server()
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
