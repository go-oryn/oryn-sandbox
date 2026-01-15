package healthcheck

import (
	"context"
	"errors"
	"sync"
)

type Probe interface {
	Name() string
	Probe(ctx context.Context) (string, error)
}

type Checker struct {
	probes []Probe
}

func NewChecker(probes ...Probe) *Checker {
	return &Checker{
		probes: probes,
	}
}

func (c *Checker) Check(ctx context.Context) (map[string]string, error) {
	success := true
	results := make(map[string]string, len(c.probes))

	var wg sync.WaitGroup
	wg.Add(len(c.probes))

	var mut sync.Mutex

	for _, probe := range c.probes {
		go func(ctx context.Context, p Probe) {
			defer wg.Done()

			res, err := p.Probe(ctx)
			mut.Lock()
			if err != nil {
				success = false
			}
			results[p.Name()] = res
			mut.Unlock()
		}(ctx, probe)
	}

	wg.Wait()

	if !success {
		return results, errors.New("healthcheck failed")
	}

	return results, nil
}
