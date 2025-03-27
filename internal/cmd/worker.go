package cmd

import (
	"ansible-runner/internal/tasks"
	"ansible-runner/internal/utils"
	"github.com/spf13/cobra"
)

var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Launch a worker instance to hanle an Ansible jobs",
	Run: func(cmd *cobra.Command, args []string) {
		utils.LoadConfig("./configs/config")
		utils.InitLogger()
		if err := tasks.LaunchServer(); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(workerCmd)
	// add a flag to set the port
	//serverCmd.Flags().StringP("port", "p", "8080", "Port to listen on")

	// bind the flag to the viper config
	//viper.BindPFlag("server.port", serverCmd.Flags().Lookup("port"))
}
