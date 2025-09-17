package cmd

import (
	"ansible-runner/internal/db"
	"ansible-runner/internal/tasks"
	"ansible-runner/internal/utils"
	"os"
	"os/signal"
	"syscall"
	"github.com/spf13/cobra"
)

// PostWorkerRun is a function that is called after the worker command is run.
func PostWorkerRun(cmd *cobra.Command, args []string) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	utils.Log.Info("Cleanup and shutting down worker...")
	tasks.Close()
	tasks.CloseClient()
	db.MongoClient.CloseService()
	utils.CleanUp()
}

var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Launch a worker instance to hanle an Ansible jobs",
	Run: func(cmd *cobra.Command, args []string) {
		utils.LoadConfig("./configs/config")
		utils.InitLogger()

		go func ()  {
			if err := tasks.LaunchServer("worker"); err != nil {
				utils.Log.Error(err.Error())
				panic(err)
			}
		}()
	},
	PostRun: PostWorkerRun,
}

func init() {
	rootCmd.AddCommand(workerCmd)
	// add a flag to set the port
	//serverCmd.Flags().StringP("port", "p", "8080", "Port to listen on")

	// bind the flag to the viper config
	//viper.BindPFlag("server.port", serverCmd.Flags().Lookup("port"))
}
