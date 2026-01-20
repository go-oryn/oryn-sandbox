package api

import (
	"github.com/go-oryn/oryn-sandbox/internal"
	"github.com/go-oryn/oryn-sandbox/pkg/healthcheck"
	"github.com/go-oryn/oryn-sandbox/pkg/httpserver"
	"github.com/go-oryn/oryn-sandbox/pkg/mcpserver"
	"github.com/spf13/cobra"
)

var ServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve API",
	Run: func(cmd *cobra.Command, args []string) {
		internal.Run(
			cmd.Context(),
			//fx.NopLogger,
			healthcheck.RunServer(),
			httpserver.RunServer(),
			mcpserver.RunStreamableHTTPServer(),
			//worker.RunWorkers(),
		)
	},
}
