package cyqldog

import (
	"log"
	"time"
)

// Scheduler represents the schedule corresponding to the rule.
type Scheduler struct {
	id   int
	rule Rule
}

// newScheduler returns an instance of Scheduler.
func newScheduler(id int, rule Rule) *Scheduler {
	return &Scheduler{
		id:   id,
		rule: rule,
	}
}

// run periodically generates monitoring tasks according to the rule.
func (s *Scheduler) run(q chan<- Rule) {
	log.Printf("scheduler(%d): start", s.id)

	// Generate trigger periodically.
	t := time.NewTicker(s.rule.Interval)
	defer t.Stop()

	// If the monitoring interval is long,
	// it will take time to check whether it is in the normal state,
	// so monitor once after startup.
	log.Printf("scheduler(%d): check on startup: %s", s.id, s.rule.Name)
	q <- s.rule

	for {
		select {
		case <-t.C:
			log.Printf("scheduler(%d): triggered: %s", s.id, s.rule.Name)
			// So as not to consume the database connection simultaneously
			// among the schedulers with different intervals,
			// we put a task in the queue and serialize the monitoring.
			// Taking into account the case of the monitoring query is slow,
			// block here without buffers to prevent duplicate monitoring tasks.
			q <- s.rule
		}
	}
}
