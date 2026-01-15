package db

import (
	"strings"

	"github.com/go-oryn/oryn-sandbox/internal"
	"github.com/go-oryn/oryn-sandbox/pkg/db"
	"github.com/spf13/cobra"
)

var MigrateCmd = &cobra.Command{
	Use:     "migrate",
	Short:   "Migrate database",
	Example: strings.Join(migrateExamples, "\n"),
	Run: func(cmd *cobra.Command, args []string) {
		internal.Run(
			cmd.Context(),
			//fx.NopLogger,
			db.RunMigrationsAndShutdown(args[0], args[1:]...),
		)
	},
}

var migrateExamples = []string{
	"  migrate up                    # migrate the DB to the most recent version available",
	"  migrate up-by-one             # migrate the DB up by 1",
	"  migrate up-to VERSION         # migrate the DB to a specific VERSION",
	"  migrate down                  # roll back the version by 1",
	"  migrate down-to VERSION       # roll back to a specific VERSION",
	"  migrate redo                  # re-run the latest migration",
	"  migrate reset                 # roll back all migrations",
	"  migrate status                # dump the migration status for the current DB",
	"  migrate version               # print the current version of the database",
	"  migrate create NAME [sql|go]  # creates new migration file with the current timestamp",
	"  migrate fix                   # apply sequential ordering to migrations",
	"  migrate validate              # check migration files without running them",
}
