package cyqldog

import (
	"database/sql"
	"log"

	"github.com/pkg/errors"
)

// DB is an implementation of DataSource.
type DB struct {
	db *sql.DB
}

// newDB returns an instance of DataSource interface.
// This function returns a error if the connection test fails.
func newDB(s DataSourceConfig) (DataSource, error) {

	// Join the options into a string of data source name.
	dataSourceName, err := s.getDataSourceName()
	if err != nil {
		return nil, err
	}

	// Open the database.
	// Note that network connection is not established at this time.
	db, err := sql.Open(s.Driver, dataSourceName)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open database")
	}

	// Connect to the database and verify its connection.
	if err = db.Ping(); err != nil {
		db.Close()

		return nil, errors.Wrapf(err, "failed to connect database")
	}

	return &DB{db: db}, nil
}

// Get queries the database to generate metrics.
func (d *DB) Get(rule Rule) (result, error) {
	r := result{}

	// Execute the SQL.
	log.Printf("db: query: %s", rule.Query)
	rows, err := d.db.Query(rule.Query)
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

// Close closes the database connection.
func (d *DB) Close() error {
	return d.db.Close()
}
