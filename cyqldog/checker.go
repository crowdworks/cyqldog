package cyqldog

import (
	"database/sql"
	"log"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/pkg/errors"
)

// Checker is a worker that executes SQLs and sends metrics.
type Checker struct {
	db     *sql.DB
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
func newChecker(db *sql.DB, statsd *statsd.Client) *Checker {
	return &Checker{
		db:     db,
		statsd: statsd,
	}
}

// run processes the monitoring task queue enqueued by the Scheduler.
func (m *Checker) run(q chan Rule) {
	log.Printf("checker: start")

	for {
		select {
		case rule := <-q:
			// dequeue the task and check.
			log.Printf("checker: check: %s", rule.Name)
			err := m.check(rule)
			if err != nil {
				// Currently, there is no way to notify errors.
				// So we simply exit the program for the time being.
				log.Fatal(err)
			}
		}
	}
}

// check gets the metrics and sends them.
func (m *Checker) check(rule Rule) error {
	result, err := m.get(rule)
	if err != nil {
		return err
	}
	return m.put(result)
}

// get queries the database to generate metrics.
func (m *Checker) get(rule Rule) (result, error) {
	r := result{}

	// Execute the SQL.
	log.Printf("checker: query: %s", rule.Query)
	rows, err := m.db.Query(rule.Query)
	if err != nil {
		return r, errors.Wrapf(err, "failed to query: %s", rule.Query)
	}
	defer rows.Close()

	// Get columns to map metric values and tags.
	cols, err := rows.Columns()
	if err != nil {
		return r, errors.Wrapf(err, "failed to get column names: %s", rule.Query)
	}

	// For each rows.
	for rows.Next() {
		// At this point, the type of each column is unknown.
		// So we temporarily store values in a slice of interface.
		row := make([]interface{}, len(cols))
		rowPtr := make([]interface{}, len(cols))
		for i := range row {
			// Scan requires a slice of pointers.
			rowPtr[i] = &row[i]
		}
		err := rows.Scan(rowPtr...)
		if err != nil {
			return r, errors.Wrapf(err, "failed to scan value: %s", rule.Query)
		}

		// Convert the row to metrics.
		metrics, err := rule.buildMetrics(row, cols)
		if err != nil {
			return r, err
		}

		// Add metrics to the result.
		r.metrics = append(r.metrics, metrics...)
	}
	return r, err
}

// put sends metrics to the dogstatsd.
func (m *Checker) put(r result) error {
	// For each metric.
	for _, metric := range r.metrics {
		log.Printf("checker: put: %s(%s) = %v\n", metric.name, metric.tags, metric.value)

		// Send a metic to the dogstatsd.
		err := m.statsd.Gauge(metric.name, metric.value, metric.tags, 1)
		if err != nil {
			return errors.Wrapf(err, "failed to gauge statsd for name = %s, value = %v, tags = %v", metric.name, metric.value, metric.tags)
		}
	}
	return nil
}
