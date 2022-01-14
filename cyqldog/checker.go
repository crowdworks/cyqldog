package cyqldog

import (
	"log"
)

// Checker is a worker that executes SQLs and sends metrics.
type Checker struct {
	ds        DataSource
	notifiers Notifiers
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
func newChecker(ds DataSource, notifiers Notifiers) *Checker {
	return &Checker{
		ds:        ds,
		notifiers: notifiers,
	}
}

// run processes the monitoring task queue enqueued by the Scheduler.
func (c *Checker) run(q <-chan Rule) {
	log.Printf("checker: start")

	for {
		rule := <-q
		log.Printf("checker: check: %s", rule.Name)

		// dequeue the task and check.
		if err := c.check(rule); err != nil {
			log.Printf("checker: failed to check: %+v", err)

			// send an error event to the notifier.
			event := newErrorEvent(err)
			if err := c.notifiers[rule.Notifier].Event(event); err != nil {
				// Sending error event was failed.
				// There is no way to notify errors, so we simply exit the program.
				log.Fatalf("failed to send error event: %+v", err)
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
	return c.notifiers[rule.Notifier].Put(result, rule)
}
