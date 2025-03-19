package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ansible-runner",
	Short: "Ansible Runner is a tool to run Ansible playbooks",
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
