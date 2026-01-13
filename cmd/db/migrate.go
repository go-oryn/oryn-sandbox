package db

import (
	"fmt"

	"github.com/spf13/cobra"
)

var MigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate database",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("migrate called, nothing will be executed for now")
	},
}
