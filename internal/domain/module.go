package domain

import (
	"github.com/go-oryn/oryn-sandbox/internal/domain/greet"
	"go.uber.org/fx"
)

const ModuleName = "domain"

var Module = fx.Module(
	ModuleName,
	// greet domain
	fx.Provide(
		greet.NewRepository,
		greet.NewService,
	),
)
