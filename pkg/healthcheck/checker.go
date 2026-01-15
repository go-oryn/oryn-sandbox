package healthcheck

import (
	"context"
	"log/slog"
	"sync"
)

const (
	Healthy   = "healthy"
	Unhealthy = "unhealthy"
)

type Probe interface {
	Name() string
	Probe(ctx context.Context) error
}

type ProbeResult struct {
	Status string `json:"status"`
}

type CheckerResult struct {
	Status       string                 `json:"status"`
	ProbeResults map[string]ProbeResult `json:"probes"`
}

func (r CheckerResult) Healthy() bool {
	return r.Status == Healthy
}

type Checker struct {
	logger *slog.Logger
	probes []Probe
}

func NewChecker(logger *slog.Logger, probes ...Probe) *Checker {
	return &Checker{
		logger: logger,
		probes: probes,
	}
}

func (c *Checker) Check(ctx context.Context) CheckerResult {
	result := CheckerResult{
		Status:       Healthy,
		ProbeResults: make(map[string]ProbeResult, len(c.probes)),
	}

	var wg sync.WaitGroup
	wg.Add(len(c.probes))

	var mut sync.Mutex

	for _, probe := range c.probes {
		go func(ctx context.Context, p Probe) {
			defer wg.Done()

			probeResult := ProbeResult{
				Status: Healthy,
			}

			err := p.Probe(ctx)

			mut.Lock()
			if err != nil {
				c.logger.ErrorContext(ctx, "health check probe failed", "probe", p.Name(), "error", err)

				result.Status = Unhealthy
				probeResult.Status = Unhealthy
			}

			result.ProbeResults[p.Name()] = probeResult
			mut.Unlock()
		}(ctx, probe)
	}

	wg.Wait()

	if result.Healthy() {
		c.logger.DebugContext(ctx, "health check success")
	}

	return result
}
