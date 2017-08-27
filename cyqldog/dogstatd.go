package cyqldog

import (
	"fmt"
	"strconv"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/pkg/errors"
)

// Dogstatsd is a configuration of the dogstatsd to connect.
type Dogstatsd struct {
	// Host is a hostname or IP address of the dogstatsd.
	Host string `yaml:"host"`
	// Port is a port number of the dogstatsd.
	Port string `yaml:"port"`
	// Namespace to prepend to all statsd calls
	Namespace string `yaml:"namespace"`
	// Tags are global tags to be added to every statsd call
	Tags []string `yaml:"tags"`
}

// newStatsd returns an instance of statsd.Client.
func newStatsd(d Dogstatsd) (*statsd.Client, error) {
	address := fmt.Sprintf("%s:%s", d.Host, d.Port)
	c, err := statsd.New(address)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open statsd: host=%s port=%s", d.Host, d.Port)
	}

	c.Namespace = d.Namespace + "."
	c.Tags = append(c.Tags, d.Tags...)
	return c, nil
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
			return metrics, errors.Wrapf(err, "failed to ParseFloat: col = %s, type = %T(%v)", vc, record[vc], record[vc])
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
