package cyqldog

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// DataSource is an interface which get metrics from.
type DataSource interface {
	Get(rule Rule) (QueryResult, error)
	Close() error
}

// Record is a map of column / value pairs representing one row.
type Record map[string]string

// QueryResult has multiple records.
type QueryResult struct {
	Records []Record
}

// DataSourceConfig is a configuration of the database to connect.
type DataSourceConfig struct {
	// Driver is type of the database to connect.
	// Currently suppoted databases are as follows:
	//  - postgres
	//  - mysql
	Driver string `yaml:"driver"`
	// Options is a map of options to connect.
	// These options are passed to sql.Open.
	// The supported options are depend on the database driver.
	Options DataSourceOptions `yaml:"options"`
}

// DataSourceOptions is a map of options to connect.
type DataSourceOptions map[string]string

// getDataSourceName returns a data source name to use for sql.Open.
func (s *DataSourceConfig) getDataSourceName() (string, error) {
	// Check database driver
	switch s.Driver {
	case "postgres":
		return s.getDataSourceNamePostgres()
	case "mysql":
		return s.getDataSourceNameMySQL()
	default:
		return "", errors.Errorf("unsupported database driver: %s", s.Driver)
	}
}

func (s *DataSourceConfig) getDataSourceNamePostgres() (string, error) {
	opts := make([]string, len(s.Options))
	for k, v := range s.Options {
		o := k + "=" + v
		opts = append(opts, o)
	}

	return strings.Join(opts[:], " "), nil
}

func (s *DataSourceConfig) getDataSourceNameMySQL() (string, error) {
	o := s.Options

	// render DSN format
	name := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		o["user"], o["password"], o["host"], o["port"], o["dbname"],
	)

	return name, nil
}
