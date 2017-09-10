package cyqldog

import (
	"strings"
	"testing"
)

func TestGetDataSourceName(t *testing.T) {
	cases := []struct {
		driver string
		ok     bool
	}{
		{
			driver: "postgres",
			ok:     true,
		},
		{
			driver: "mysql",
			ok:     true,
		},
		{
			driver: "unknown",
			ok:     false,
		},
	}

	for _, tc := range cases {
		s := &DataSourceConfig{
			Driver: tc.driver,
		}

		_, err := s.getDataSourceName()

		if tc.ok && err != nil {
			t.Errorf("getDataSourceName() with driver = %s returns unexpected error: %+v", tc.driver, err)
		}

		if !tc.ok && err == nil {
			t.Errorf("expected getDataSourceName() with driver = %s returns error, but err == nil", tc.driver)
		}
	}

}

func TestGetDataSourceNamePostgres(t *testing.T) {
	s := DataSourceConfig{
		Driver: "postgres",
		Options: DataSourceOptions{
			"host":     "db.example.com",
			"port":     "5432",
			"user":     "cyqldog",
			"password": "secret",
			"dbname":   "cyqldogdb",
			"sslmode":  "disable",
		},
	}

	out, _ := s.getDataSourceNamePostgres()

	// A map interation order is random, so the output is unstable.
	// We only check if the key=value is included in the output.
	for k, v := range s.Options {
		if ok := strings.Contains(out, k+"="+v); !ok {
			t.Errorf("getDataSourceNamePostgres() does not contains %s=%s", k, v)
		}
	}

}

func TestGetDataSourceNameMySQL(t *testing.T) {
	cases := []struct {
		options DataSourceOptions
		out     string
	}{
		{
			options: DataSourceOptions{
				"host":     "db.example.com",
				"port":     "3306",
				"user":     "cyqldog",
				"password": "secret",
				"dbname":   "cyqldogdb",
			},
			out: "cyqldog:secret@tcp(db.example.com:3306)/cyqldogdb",
		},
	}

	for _, tc := range cases {
		s := DataSourceConfig{
			Driver:  "mysql",
			Options: tc.options,
		}

		got, _ := s.getDataSourceNameMySQL()

		if got != tc.out {
			t.Errorf("getDataSourceNameMySQL() with options = %v returns %s, but want = %s", tc.options, got, tc.out)
		}

	}

}
