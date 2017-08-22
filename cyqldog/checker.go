package cyqldog

import (
	"log"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/pkg/errors"
)

// Checker is a worker that executes SQLs and sends metrics.
type Checker struct {
	ds     DataSource
	statsd *statsd.Client
}

// metric represents a measured value.
type metric struct {
	name  string
	value float64
	tags  []string
}

// result represents monitoring results.
type result struct {
	metrics []metric
}

// newChecker returns an instance of Checker.
func newChecker(ds DataSource, statsd *statsd.Client) *Checker {
	return &Checker{
		ds:     ds,
		statsd: statsd,
	}
}

// run processes the monitoring task queue enqueued by the Scheduler.
func (c *Checker) run(q chan Rule) {
	log.Printf("checker: start")

	for {
		select {
		case rule := <-q:
			// dequeue the task and check.
			log.Printf("checker: check: %s", rule.Name)
			err := c.check(rule)
			if err != nil {
				// Currently, there is no way to notify errors.
				// So we simply exit the program for the time being.
				log.Fatal(err)
			}
		}
	}
}

// check gets the metrics and sends them.
func (c *Checker) check(rule Rule) error {
	result, err := c.ds.Get(rule)
	if err != nil {
		return err
	}
	return c.put(result)
}

// put sends metrics to the dogstatsd.
func (c *Checker) put(r result) error {
	// For each metric.
	for _, metric := range r.metrics {
		log.Printf("checker: put: %s(%s) = %v\n", metric.name, metric.tags, metric.value)

		// Send a metic to the dogstatsd.
		err := c.statsd.Gauge(metric.name, metric.value, metric.tags, 1)
		if err != nil {
			return errors.Wrapf(err, "failed to gauge statsd for name = %s, value = %v, tags = %v", metric.name, metric.value, metric.tags)
		}
	}
	return nil
}
