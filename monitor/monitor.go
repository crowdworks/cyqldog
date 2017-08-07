package monitor

import (
	"database/sql"
	"log"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/pkg/errors"
)

type Monitor struct {
	db     *sql.DB
	statsd *statsd.Client
}

type metric struct {
	name  string
	value float64
	tags  []string
}

type result struct {
	metrics []metric
}

func NewMonitor(db *sql.DB, statsd *statsd.Client) *Monitor {
	return &Monitor{
		db:     db,
		statsd: statsd,
	}
}

func (m *Monitor) Run(q chan Rule) {
	log.Printf("monitor: start")

	for {
		select {
		case rule := <-q:
			log.Printf("monitor: check: %s", rule.Name)
			err := m.check(rule)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func (m *Monitor) check(rule Rule) error {
	result, err := m.get(rule)
	if err != nil {
		return err
	}
	return m.put(result)
}

func (m *Monitor) get(rule Rule) (result, error) {
	r := result{}

	log.Printf("monitor: query: %s", rule.Query)
	rows, err := m.db.Query(rule.Query)
	if err != nil {
		return r, errors.Wrapf(err, "failed to query: %s", rule.Query)
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return r, errors.Wrapf(err, "failed to get column names: %s", rule.Query)
	}

	for rows.Next() {
		row := make([]interface{}, len(cols))
		rowPtr := make([]interface{}, len(cols))
		for i := range row {
			rowPtr[i] = &row[i]
		}

		err := rows.Scan(rowPtr...)
		if err != nil {
			return r, errors.Wrapf(err, "failed to scan value: %s", rule.Query)
		}

		metrics, err := rule.buildMetrics(row, cols)
		if err != nil {
			return r, err
		}

		r.metrics = append(r.metrics, metrics...)
	}
	return r, err
}

func (m *Monitor) put(r result) error {
	for _, metric := range r.metrics {
		log.Printf("monitor: put: %s(%s) = %v\n", metric.name, metric.tags, metric.value)
		err := m.statsd.Gauge(metric.name, metric.value, metric.tags, 1)
		if err != nil {
			return errors.Wrapf(err, "failed to gauge statsd for name = %s, value = %v, tags = %v", metric.name, metric.value, metric.tags)
		}
	}
	return nil
}
