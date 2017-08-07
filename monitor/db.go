package monitor

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
)

type DataSource struct {
	Driver   string `yaml:"driver"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Dbname   string `yaml:"dbname"`
	Sslmode  string `yaml:"sslmode"`
}

func NewDB(s DataSource) (*sql.DB, error) {
	dataSourceName := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		s.Host, s.Port, s.User, s.Password, s.Dbname, s.Sslmode,
	)
	db, err := sql.Open(s.Driver, dataSourceName)
	if err != nil {
		return nil, errors.Wrapf(err,
			"failed to open database: driver=%s host=%s port=%s user=%s dbname=%s sslmode=%s",
			s.Driver, s.Host, s.Port, s.User, s.Dbname, s.Sslmode,
		)
	}

	if err = db.Ping(); err != nil {
		db.Close()

		return nil, errors.Wrapf(err,
			"failed to connect database: driver=%s host=%s port=%s user=%s dbname=%s sslmode=%s",
			s.Driver, s.Host, s.Port, s.User, s.Dbname, s.Sslmode,
		)
	}
	return db, nil
}
