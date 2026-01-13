package core

import (
	"github.com/go-oryn/oryn-sandbox/pkg/config"
	"github.com/go-oryn/oryn-sandbox/pkg/otel"
	"go.uber.org/fx"
)

const ModuleName = "core"

var Module = fx.Module(
	ModuleName,
	// sub modules
	config.Module,
	otel.Module,
	// configurations
	//ConfigureFx(),
	ConfigureOTel(),
)
