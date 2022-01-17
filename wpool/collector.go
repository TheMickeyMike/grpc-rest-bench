package wpool

import (
	"context"
	"fmt"
)

type Collector struct {
	resultsQueue <-chan Result
	results      map[string]int
	done         chan struct{}
}

func NewCollector(resultsQueue <-chan Result) *Collector {
	return &Collector{
		resultsQueue: resultsQueue,
		results:      make(map[string]int),
		done:         make(chan struct{}),
	}
}

func (c *Collector) Run(ctx context.Context) {
	defer close(c.done)
	for {
		select {
		case res, ok := <-c.resultsQueue:
			if !ok {
				return
			}
			if res.Err != nil {
				c.results["error"] += 1
			} else {
				c.results[res.Value] += 1
			}
		case <-ctx.Done():
			fmt.Printf("cancelled collector. Error detail: %v\n", ctx.Err())
			return
		}
	}
}

func (c *Collector) GenerateReport() map[string]int {
	<-c.done // wait until collector end runnig
	return c.results
}
