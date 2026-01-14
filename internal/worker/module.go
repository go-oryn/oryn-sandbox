package worker

import (
	"github.com/go-oryn/oryn-sandbox/internal/worker/greet"
	"github.com/go-oryn/oryn-sandbox/pkg/worker"
	"go.uber.org/fx"
)

const ModuleName = "worker"

var Module = fx.Module(
	ModuleName,
	// greet
	worker.AsWorker(greet.NewGreetWorker),
)
