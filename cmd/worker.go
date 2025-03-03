package cmd

import (
	"ansible-runner/configs"
	"ansible-runner/routers"
	"ansible-runner/utils"
	"github.com/spf13/cobra"
)

var workerCmd = &cobra.Command{
	Use:   "server",
	Short: "Launch a worker instance to hanle an Ansible jobs",
	Run: func(cmd *cobra.Command, args []string) {
		utils.InitLogger("worker.log")
		configs.LoadConfig("./")
		routers.NewWorker()
	},
}

func init() {
	rootCmd.AddCommand(workerCmd)
	// add a flag to set the port
	//serverCmd.Flags().StringP("port", "p", "8080", "Port to listen on")

	// bind the flag to the viper config
	//viper.BindPFlag("server.port", serverCmd.Flags().Lookup("port"))
}
