package cyqldog

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
)

// DataSource is a configuration of the database to connect.
type DataSource struct {
	// Driver is type of the database to connect.
	// Currently suppoted database is postgres.
	Driver string `yaml:"driver"`
	// Host is hostname or IP address of the database.
	Host string `yaml:"host"`
	// Port is a port number of the database.
	Port string `yaml:"port"`
	// User is a username of the database.
	User string `yaml:"user"`
	// Password is a password of the database.
	Password string `yaml:"password"`
	// Dbname is a name of the database.
	Dbname string `yaml:"dbname"`
	// Sslmode is an SSL connection option.
	// Suppoted values are `require`, `verify-ca`, `verify-full`, and `disable`.
	Sslmode string `yaml:"sslmode"`
}

// NewDB returns an instance of sql.DB.
// This function returns a error if the connection test fails.
func NewDB(s DataSource) (*sql.DB, error) {
	dataSourceName := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		s.Host, s.Port, s.User, s.Password, s.Dbname, s.Sslmode,
	)

	// Open the database.
	// Note that network connection is not established at this time.
	db, err := sql.Open(s.Driver, dataSourceName)
	if err != nil {
		return nil, errors.Wrapf(err,
			"failed to open database: driver=%s host=%s port=%s user=%s dbname=%s sslmode=%s",
			s.Driver, s.Host, s.Port, s.User, s.Dbname, s.Sslmode,
		)
	}

	// Connect to the database and verify its connection.
	if err = db.Ping(); err != nil {
		db.Close()

		return nil, errors.Wrapf(err,
			"failed to connect database: driver=%s host=%s port=%s user=%s dbname=%s sslmode=%s",
			s.Driver, s.Host, s.Port, s.User, s.Dbname, s.Sslmode,
		)
	}
	return db, nil
}
