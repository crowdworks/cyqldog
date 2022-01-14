package cyqldog

import (
	"strings"

	"golang.org/x/xerrors"

	"github.com/go-sql-driver/mysql"
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
		return "", xerrors.Errorf("unsupported database driver: %s", s.Driver)
	}
}

func (s *DataSourceConfig) getDataSourceNamePostgres() (string, error) {
	opts := []string{}
	for k, v := range s.Options {
		o := k + "=" + v
		opts = append(opts, o)
	}

	return strings.Join(opts[:], " "), nil
}

func (s *DataSourceConfig) getDataSourceNameMySQL() (string, error) {
	o := s.Options

	port := "3306"
	if len(o["port"]) > 0 {
		port = o["port"]
	}

	c := mysql.NewConfig()

	c.User = o["user"]
	c.Passwd = o["password"]
	c.Net = "tcp"
	c.Addr = o["host"] + ":" + port
	c.DBName = o["dbname"]

	// delete basic options.
	delete(o, "user")
	delete(o, "password")
	delete(o, "host")
	delete(o, "port")
	delete(o, "dbname")

	// set other connection params.
	c.Params = o

	// render DSN format
	return c.FormatDSN(), nil
}
