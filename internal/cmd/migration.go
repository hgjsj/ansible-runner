package cmd

// add subcommands migration
import (
	"ansible-runner/internal/db"
	"ansible-runner/internal/utils"
	"fmt"
	"github.com/spf13/cobra"
)

var migrationCmd = &cobra.Command{
	Use:   "migration",
	Short: "Run database migrations",
	Run: func(cmd *cobra.Command, args []string) {
		// Call the migration function here
		runMigrations()
	},
}

func init() {
	rootCmd.AddCommand(migrationCmd)
}

func runMigrations() {
	// Implement your migration logic here
	utils.LoadConfig("./configs/config")
	utils.InitLogger()
	if err := db.GetDBClient().Migration(); err != nil {
		fmt.Println("Migration failed:", err.Error())
	}
	defer db.GetDBClient().CloseService()

}
