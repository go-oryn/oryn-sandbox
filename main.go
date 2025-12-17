package main

import (
	"fmt"

	"github.com/go-oryn/oryn-sandbox/configs"
	"github.com/go-oryn/oryn-sandbox/pkg/config"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		config.ConfigModule,
		config.AsConfigOptions(config.WithEmbedFS(configs.ConfigFS)),
		fx.Invoke(func(cfg *config.Config, shutdown fx.Shutdowner) error {
			fmt.Printf("App name: %s\n", cfg.GetString("app.name"))

			return shutdown.Shutdown()
		}),
	).Run()
}
