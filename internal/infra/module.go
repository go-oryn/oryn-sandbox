package worker

import (
	"github.com/go-oryn/oryn-sandbox/pkg/db"
	"github.com/go-oryn/oryn-sandbox/pkg/healthcheck"
	"github.com/go-oryn/oryn-sandbox/pkg/worker"
	"go.uber.org/fx"
)

const ModuleName = "infra"

var Module = fx.Module(
	ModuleName,
	// probes
	healthcheck.AsProbe(db.NewDBProbe),
	healthcheck.AsProbe(worker.NewWorkersProbe),
)
