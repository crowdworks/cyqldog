package cyqldog

import (
	"fmt"
	"golang.org/x/xerrors"
	"log"
	"strconv"

	"github.com/DataDog/datadog-go/statsd"
)

// DogstatsdConfig is a configuration of the dogstatsd to connect.
type DogstatsdConfig struct {
	// Host is a hostname or IP address of the dogstatsd.
	Host string `yaml:"host"`
	// Port is a port number of the dogstatsd.
	Port string `yaml:"port"`
	// Namespace to prepend to all statsd calls
	Namespace string `yaml:"namespace"`
	// Tags are global tags to be added to every statsd call
	Tags []string `yaml:"tags"`
}

// statsdClient is an interface of statsd.Client.
// We make a layer of abstraction for testing.
type statsdClient interface {
	Gauge(name string, value float64, tags []string, rate float64) error
	Event(e *statsd.Event) error
}

// Dogstatsd is a configuration of the dogstatsd to connect.
type Dogstatsd struct {
	client statsdClient
}

// newDogstatsd returns an instance of Notifier interface.
func newDogstatsd(d DogstatsdConfig) (Notifier, error) {
	address := fmt.Sprintf("%s:%s", d.Host, d.Port)
	c, err := statsd.New(address)
	if err != nil {
		return nil, xerrors.Errorf("failed to open statsd: host=%s port=%s: %w", d.Host, d.Port, err)
	}

	c.Namespace = d.Namespace + "."
	c.Tags = append(c.Tags, d.Tags...)
	return &Dogstatsd{client: c}, nil
}

// Event send an event to the dogstatsd.
func (d *Dogstatsd) Event(e *Event) error {
	se := &statsd.Event{
		Title:          e.Title,
		Text:           e.Text,
		AggregationKey: "cyqldog",
		Tags:           e.Tags,
	}

	// Unfortunately statsd.eventAlertType is an unexported type.
	// However the constants are exported,
	// so we convert string to statsd.eventAlertType by ourselves.
	switch e.Level {
	case "", "info": // default is info
		se.AlertType = statsd.Info
	case "error":
		se.AlertType = statsd.Error
	case "warning":
		se.AlertType = statsd.Warning
	case "success":
		se.AlertType = statsd.Success
	default:
		return xerrors.Errorf("unknown event level: %+v", e)
	}

	return d.client.Event(se)
}

// Put sends metrics to the dogstatsd.
func (d *Dogstatsd) Put(qr QueryResult, rule Rule) error {
	// convert to the query result to metrics.
	metrics, err := buildMetricsForQueryResult(qr, rule)
	if err != nil {
		return err
	}

	// For each metric.
	for _, metric := range metrics {
		log.Printf("checker: put: %s(%s) = %v\n", metric.name, metric.tags, metric.value)

		// Send a metic to the dogstatsd.
		err := d.client.Gauge(metric.name, metric.value, metric.tags, 1)
		if err != nil {
			return xerrors.Errorf("failed to gauge statsd for name = %s, value = %v, tags = %v: %w", metric.name, metric.value, metric.tags, err)
		}
	}
	return nil
}

// buildMetricsForRecord returns a metrics from the query result.
func buildMetricsForQueryResult(qr QueryResult, rule Rule) ([]metric, error) {
	metrics := []metric{}

	for _, record := range qr.Records {
		// convert record to metrics.
		ms, err := buildMetricsForRecord(record, rule.Name, rule.ValueCols, rule.TagCols)
		if err != nil {
			return metrics, err
		}

		metrics = append(metrics, ms...)
	}

	return metrics, nil
}

// buildMetricsForRecord returns a metrics from the record.
func buildMetricsForRecord(record Record, prefix string, valueCols []string, tagCols []string) ([]metric, error) {
	metrics := []metric{}

	for _, vc := range valueCols {
		// The value of record is a string for general purpose,
		// so we parse and convert it to float64 here.
		value, err := strconv.ParseFloat(record[vc], 64)
		if err != nil {
			return metrics, xerrors.Errorf("failed to ParseFloat: col = %s, type = %T(%v): %w", vc, record[vc], record[vc], err)
		}

		// build metric.
		m := metric{
			name:  prefix + "." + vc, // prefix is Rule.Name
			value: value,
			tags:  buildTags(record, tagCols),
		}
		metrics = append(metrics, m)
	}

	return metrics, nil
}

// buildTags returns a slice of tags from the record and column names to use for tag.
func buildTags(record Record, tagCols []string) []string {
	tags := []string{}

	for _, tc := range tagCols {
		// tags are formatted as column name:value.
		tags = append(tags, tc+":"+record[tc])
	}

	return tags
}
