package cyqldog

import (
	"database/sql/driver"
	"reflect"
	"regexp"
	"testing"
	"time"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestDBGet(t *testing.T) {
	cases := []struct {
		in       Rule
		mockCols []string
		mockRows [][]driver.Value
		out      QueryResult
	}{
		{
			in: Rule{
				Name:      "test1",
				Interval:  (5 * time.Second),
				Query:     "SELECT COUNT(*) AS count FROM table1",
				Notifier:  "dogstatsd",
				ValueCols: []string{"count"},
				TagCols:   []string{},
			},
			mockCols: []string{"count"},
			mockRows: [][]driver.Value{{int64(3)}},
			out: QueryResult{
				Records: []Record{
					{"count": "3"},
				},
			},
		},
		{
			in: Rule{
				Name:      "test2",
				Interval:  (10 * time.Second),
				Query:     "SELECT tag1, val1, tag2, val2 FROM table1",
				Notifier:  "dogstatsd",
				ValueCols: []string{"tag1", "tag2"},
				TagCols:   []string{"val1", "val2"},
			},
			mockCols: []string{"tag1", "val1", "tag2", "val2"},
			mockRows: [][]driver.Value{
				{"hoge1", int64(1), "fuga1", float64(0.1)},
				{"hoge1", int64(2), "fuga2", float64(0.2)},
				{"hoge3", int64(3), "fuga3", float64(0.3)},
			},
			out: QueryResult{
				Records: []Record{
					{"tag1": "hoge1", "val1": "1", "tag2": "fuga1", "val2": "0.100000"},
					{"tag1": "hoge1", "val1": "2", "tag2": "fuga2", "val2": "0.200000"},
					{"tag1": "hoge3", "val1": "3", "tag2": "fuga3", "val2": "0.300000"},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.in.Name, func(t *testing.T) {
			// Create mockDB
			mockDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to open mock database: %v", err)
			}
			defer mockDB.Close()

			d := &DB{db: mockDB}

			mockRows := sqlmock.NewRows(tc.mockCols)
			for _, row := range tc.mockRows {
				mockRows.AddRow(row...)
			}
			mock.ExpectQuery(regexp.QuoteMeta(tc.in.Query)).WillReturnRows(mockRows)

			got, err := d.Get(tc.in)

			if err != nil {
				t.Errorf("DB.Get(%v) returns unexpected err = %+v", tc.in, err)
			}

			if !reflect.DeepEqual(got, tc.out) {
				t.Errorf("DB.Get(%v) = %v; want = %v", tc.in, got, tc.out)
			}

		})
	}
}
