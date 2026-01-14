package api

import (
	"github.com/go-oryn/oryn-sandbox/internal"
	"github.com/go-oryn/oryn-sandbox/pkg/httpserver"
	"github.com/go-oryn/oryn-sandbox/pkg/worker"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

var ServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve API",
	Run: func(cmd *cobra.Command, args []string) {
		internal.Run(
			cmd.Context(),
			fx.NopLogger,
			httpserver.RunServer(),
			worker.RunWorkers(),
		)
	},
}
