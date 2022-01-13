package cyqldog

import (
	"database/sql"
	"fmt"
	"golang.org/x/xerrors"
	"log"
)

// DB is an implementation of DataSource.
type DB struct {
	db *sql.DB
}

// newDB returns an instance of DataSource interface.
// This function returns a error if the connection test fails.
func newDB(c DataSourceConfig) (DataSource, error) {

	// Join the options into a string of data source name.
	dataSourceName, err := c.getDataSourceName()
	if err != nil {
		return nil, err
	}

	// Open the database.
	// Note that network connection is not established at this time.
	db, err := sql.Open(c.Driver, dataSourceName)
	if err != nil {
		return nil, xerrors.Errorf("failed to open database: %w", err)
	}

	// Connect to the database and verify its connection.
	if err = db.Ping(); err != nil {
		db.Close()

		return nil, xerrors.Errorf("failed to connect database: %w", err)
	}

	return &DB{db: db}, nil
}

// Get queries the database to generate metrics.
func (d *DB) Get(rule Rule) (QueryResult, error) {
	qr := QueryResult{}

	// Execute the SQL.
	log.Printf("db: query: %s", rule.Query)
	rows, err := d.db.Query(rule.Query)
	if err != nil {
		return qr, xerrors.Errorf("failed to query: %s: %w", rule.Query, err)
	}
	defer rows.Close()

	// Get columns to map metric values and tags.
	cols, err := rows.Columns()
	if err != nil {
		return qr, xerrors.Errorf("failed to get column names: %s: %w", rule.Query, err)
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
			return qr, xerrors.Errorf("failed to scan value: %s: %w", rule.Query, err)
		}

		// Convert the row to record.
		record, err := buildRecord(row, cols)
		if err != nil {
			return qr, err
		}

		// Add the record to the result.
		qr.Records = append(qr.Records, record)
	}

	return qr, nil
}

// buildRecord converts a row to record.
// The record stores all data as a string to generalize how to handle data at notifications.
func buildRecord(row []interface{}, cols []string) (Record, error) {
	record := make(Record, len(cols))

	for i, c := range row {
		s, err := convertToString(c)
		if err != nil {
			return record, xerrors.Errorf("faied to convertToString: col = %s, type = %T(%v): %w", cols[i], c, c, err)
		}

		record[cols[i]] = s
	}

	return record, nil
}

// convertToString casts interface to string.
func convertToString(i interface{}) (string, error) {
	switch s := i.(type) {
	case string:
		return s, nil
	case []uint8:
		return string(s), nil
	case int64:
		return fmt.Sprintf("%d", s), nil
	case float64:
		// Suppress the trailing zeros.
		return fmt.Sprintf("%v", s), nil
	default:
		return "", xerrors.New("failed to cast interface to string")
	}
}

// Close closes the database connection.
func (d *DB) Close() error {
	return d.db.Close()
}
