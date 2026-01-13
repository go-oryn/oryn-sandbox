package config

import "embed"

type Options struct {
	env     string
	embedFS *embed.FS
	values  map[string]any
}

type Option func(*Options) error

func WithEnvironment(env string) Option {
	return func(o *Options) error {
		o.env = env

		return nil
	}
}

func WithEmbedFS(fs embed.FS) Option {
	return func(o *Options) error {
		o.embedFS = &fs

		return nil
	}
}

func WithValues(values map[string]any) Option {
	return func(o *Options) error {
		o.values = values

		return nil
	}
}
