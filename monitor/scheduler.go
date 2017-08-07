package monitor

import (
	"log"
	"time"
)

type Scheduler struct {
	id   int
	rule Rule
}

func NewScheduler(id int, rule Rule) *Scheduler {
	return &Scheduler{
		id:   id,
		rule: rule,
	}
}

func (s *Scheduler) Run(q chan Rule) {
	log.Printf("scheduler(%d): start", s.id)

	t := time.NewTicker(s.rule.Interval)
	defer t.Stop()

	log.Printf("scheduler(%d): check on startup: %s", s.id, s.rule.Name)
	q <- s.rule

	for {
		select {
		case <-t.C:
			log.Printf("scheduler(%d): triggered: %s", s.id, s.rule.Name)
			q <- s.rule
		}
	}
}
