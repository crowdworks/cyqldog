package cyqldog

import (
	"reflect"
	"testing"
	"time"
)

// mockStatsdClient is a mock of statsd.Client.
type mockStatsdClient struct {
	logs []mockStatsdLog
}

// mockStatsdLog records the method name and arguments of the mockStatsdClient's API call for testing.
type mockStatsdLog struct {
	method string
	name   string
	value  interface{}
	tags   []string
	rate   float64
}

// Gauge implements an interface of statsdClient for testing.
// It only records API calls and does not call real dogstatsd.
func (c *mockStatsdClient) Gauge(name string, value float64, tags []string, rate float64) error {
	log := mockStatsdLog{
		method: "gauge",
		name:   name,
		value:  value,
		tags:   tags,
		rate:   rate,
	}
	c.logs = append(c.logs, log)
	return nil
}

func newMockStatsdClient() *mockStatsdClient {
	return &mockStatsdClient{}
}

func newMockDogstatsd(c *mockStatsdClient) Notifier {
	return &Dogstatsd{client: c}
}

func TestDogstatsdPut(t *testing.T) {
	cases := []struct {
		qr   QueryResult
		rule Rule
		out  []mockStatsdLog
	}{
		{
			qr: QueryResult{
				Records: []Record{
					{"count": "3"},
				},
			},
			rule: Rule{
				Name:      "test1",
				Interval:  (5 * time.Second),
				Query:     "SELECT COUNT(*) AS count FROM table1",
				Notifier:  "dogstatsd",
				ValueCols: []string{"count"},
				TagCols:   []string{},
			},
			out: []mockStatsdLog{
				{method: "gauge", name: "test1.count", value: float64(3), tags: []string{}, rate: 1},
			},
		},
		{
			qr: QueryResult{
				Records: []Record{
					{"tag1": "hoge1", "val1": "1", "tag2": "fuga1", "val2": "0.1"},
					{"tag1": "hoge1", "val1": "2", "tag2": "fuga2", "val2": "0.2"},
					{"tag1": "hoge3", "val1": "3", "tag2": "fuga3", "val2": "0.3"},
				},
			},
			rule: Rule{
				Name:      "test2",
				Interval:  (10 * time.Second),
				Query:     "SELECT tag1, val1, tag2, val2 FROM table1",
				Notifier:  "dogstatsd",
				ValueCols: []string{"val1", "val2"},
				TagCols:   []string{"tag1", "tag2"},
			},
			out: []mockStatsdLog{
				{method: "gauge", name: "test2.val1", value: float64(1), tags: []string{"tag1:hoge1", "tag2:fuga1"}, rate: 1},
				{method: "gauge", name: "test2.val2", value: float64(0.1), tags: []string{"tag1:hoge1", "tag2:fuga1"}, rate: 1},
				{method: "gauge", name: "test2.val1", value: float64(2), tags: []string{"tag1:hoge1", "tag2:fuga2"}, rate: 1},
				{method: "gauge", name: "test2.val2", value: float64(0.2), tags: []string{"tag1:hoge1", "tag2:fuga2"}, rate: 1},
				{method: "gauge", name: "test2.val1", value: float64(3), tags: []string{"tag1:hoge3", "tag2:fuga3"}, rate: 1},
				{method: "gauge", name: "test2.val2", value: float64(0.3), tags: []string{"tag1:hoge3", "tag2:fuga3"}, rate: 1},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.rule.Name, func(t *testing.T) {
			c := newMockStatsdClient()
			d := newMockDogstatsd(c)
			if err := d.Put(tc.qr, tc.rule); err != nil {
				t.Errorf("Dogstatsd.Put(%+v, %+v) retruns unexpected err = %+v", tc.qr, tc.rule, err)
			}

			if !reflect.DeepEqual(c.logs, tc.out) {
				t.Errorf("Dogstatsd.Put(%+v, %+v)\n got = %+v,\nwant = %+v", tc.qr, tc.rule, c.logs, tc.out)
			}
		})
	}

}
