package db

import "embed"

type MigratorOptions struct {
	embedFS *embed.FS
}

type MigratorOption func(*MigratorOptions)

func WithMigrationsEmbedFS(fs embed.FS) MigratorOption {
	return func(o *MigratorOptions) {
		o.embedFS = &fs
	}
}
