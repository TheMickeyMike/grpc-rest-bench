package wpool

import (
	"context"
	"fmt"
	"testing"
)

type Collector struct {
	resultsQueue <-chan Result
	results      map[string]int64
	done         chan struct{}
}

func NewCollector(resultsQueue <-chan Result) *Collector {
	return &Collector{
		resultsQueue: resultsQueue,
		results: map[string]int64{
			"retries": 0,
			"error":   0,
		},
		done: make(chan struct{}),
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
			if res.Retries > 0 {
				c.results["retries"] += res.Retries
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

func (c *Collector) GenerateReport(b *testing.B) string {
	<-c.done // wait until collector ends collecting results
	var (
		summary int64
		result  string
	)
	for k, v := range c.results {
		if k != "retries" {
			summary += v
		}
		result += fmt.Sprintf("%s: %d\n", k, v)
		b.ReportMetric(float64(v), k)
	}
	result += fmt.Sprintf("summary: %d\n", summary)
	return result
}
