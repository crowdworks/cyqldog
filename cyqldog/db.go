package cyqldog

import (
	"database/sql"
	"strings"

	"github.com/pkg/errors"
)

// DataSource is a configuration of the database to connect.
type DataSource struct {
	// Driver is type of the database to connect.
	// Currently suppoted database is postgres.
	Driver string `yaml:"driver"`
	// Options is a map of options to connect.
	// These options are passed to sql.Open.
	// The supported options are depend on the database driver.
	Options DataSourceOptions `yaml:"options"`
}

// DataSourceOptions is a map of options to connect.
type DataSourceOptions map[string]string

// newDB returns an instance of sql.DB.
// This function returns a error if the connection test fails.
func newDB(s DataSource) (*sql.DB, error) {

	// Join the options into a string of data source name.
	dataSourceName := s.getDataSourceName()

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
	return db, nil
}

// getDataSourceName returns a data source name to use for sql.Open.
func (s *DataSource) getDataSourceName() string {
	opts := make([]string, len(s.Options))
	for k, v := range s.Options {
		o := k + "=" + v
		opts = append(opts, o)
	}

	return strings.Join(opts[:], " ")
}
