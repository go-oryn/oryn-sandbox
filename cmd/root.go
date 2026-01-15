package cmd

import (
	"fmt"
	"os"

	"github.com/go-oryn/oryn-sandbox/cmd/api"
	"github.com/go-oryn/oryn-sandbox/cmd/db"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(api.ServeCmd)
	RootCmd.AddCommand(db.MigrateCmd)
	RootCmd.AddCommand(db.SeedCmd)

}

var RootCmd = &cobra.Command{
	Use: "app",
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
}
