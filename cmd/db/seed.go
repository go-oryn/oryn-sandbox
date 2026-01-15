package db

import (
	"github.com/go-oryn/oryn-sandbox/internal"
	"github.com/go-oryn/oryn-sandbox/pkg/db"
	"github.com/spf13/cobra"
)

var SeedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed database",
	Run: func(cmd *cobra.Command, args []string) {
		internal.Run(
			cmd.Context(),
			//fx.NopLogger,
			db.RunSeedsAndShutdown(args...),
		)
	},
}

var seedExamples = []string{
	"  seed           # run all seeds",
	"  seed foo bar   # run foo and bar seeds only",
}
