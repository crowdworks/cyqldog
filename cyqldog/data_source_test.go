package cyqldog

import (
	"reflect"
	"sort"
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
		base    string
		params  string
	}{
		{
			options: DataSourceOptions{
				"host":     "base.db.example.com",
				"port":     "3306",
				"user":     "cyqldog",
				"password": "secret",
				"dbname":   "cyqldogdb",
			},
			base:   "cyqldog:secret@tcp(base.db.example.com:3306)/cyqldogdb",
			params: "",
		},
		{
			options: DataSourceOptions{
				"host":     "short.db.example.com",
				"user":     "cyqldog",
				"password": "secret",
			},
			base:   "cyqldog:secret@tcp(short.db.example.com:3306)/",
			params: "",
		},
		{
			options: DataSourceOptions{
				"host":      "params.db.example.com",
				"port":      "3306",
				"user":      "cyqldog",
				"password":  "secret",
				"dbname":    "cyqldogdb",
				"charset":   "utf8",
				"collation": "utf8_general_ci",
			},
			base:   "cyqldog:secret@tcp(params.db.example.com:3306)/cyqldogdb",
			params: "charset=utf8&collation=utf8_general_ci",
		},
	}

	for _, tc := range cases {
		s := DataSourceConfig{
			Driver:  "mysql",
			Options: tc.options,
		}

		got, _ := s.getDataSourceNameMySQL()
		splittedDNS := strings.SplitN(got, "?", 2)
		base := splittedDNS[0]

		if base != tc.base {
			t.Errorf("base of getDataSourceNameMySQL() with options = %v returns %s, but want = %s", tc.options, base, tc.base)
		}

		// want no params.
		if len(tc.params) == 0 {
			continue
		}

		// got no params.
		if len(splittedDNS) < 2 {
			t.Errorf("getDataSourceNameMySQL() with options = %v returns no params, but want = %s", tc.options, tc.params)
			continue
		}

		// compare params.
		params := splittedDNS[1]
		if !splittedStringEqual(params, tc.params, "&") {
			t.Errorf("params of getDataSourceNameMySQL() with options = %v returns %s, but want = %s", tc.options, params, tc.params)
		}

	}

}

// splittedStringEqual compare whether or not strings splitted by separater are
// the same regardless of the order.
func splittedStringEqual(s1 string, s2 string, sep string) bool {
	ss1 := strings.Split(s1, sep)
	ss2 := strings.Split(s2, sep)
	sort.Strings(ss1)
	sort.Strings(ss2)
	return reflect.DeepEqual(ss1, ss2)
}
