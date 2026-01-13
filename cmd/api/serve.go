package api

import (
	"github.com/go-oryn/oryn-sandbox/internal"
	"github.com/go-oryn/oryn-sandbox/pkg/httpserver"
	"github.com/spf13/cobra"
)

var ServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve API",
	Run: func(cmd *cobra.Command, args []string) {
		internal.Run(
			cmd.Context(),
			httpserver.RunServer(),
		)
	},
}
