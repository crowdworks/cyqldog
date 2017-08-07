package monitor

import (
	"time"

	"github.com/pkg/errors"
)

// Rule represents monitoring conditions.
type Rule struct {
	// Name of the rule.
	// The metic name sent to the dogstatsd is:
	//   Dogstatsd.Namespace + Rule.Name + Rule.ValueCols[*}
	Name string `yaml:"name"`
	// Interval of the monitoring.
	Interval time.Duration `yaml:"interval"`
	// Query to the database.
	Query string `yaml:"query"`
	// ValueCols is a list of names of the columns used as metric values.
	ValueCols []string `yaml:"value_cols"`
	// TagCols is a list of names of the columns used as metric tags.
	TagCols []string `yaml:"tag_cols"`
}

// buildMetrics converts a row to metrics.
func (r *Rule) buildMetrics(row []interface{}, cols []string) ([]metric, error) {
	metrics := []metric{}
	values := []float64{}
	tags := []string{}

	// For each column.
	for i, c := range row {
		// If the column is used as metric values.
		if contains(r.ValueCols, cols[i]) {
			switch f := c.(type) {
			case float64:
				values = append(values, f)
			case int64: // integer
				values = append(values, float64(f))
			default:
				return metrics, errors.Errorf("failed to cast from interface to float64: col = %s, type = %T(%v)", cols[i], c, c)
			}
		}

		// If the column is used as metric tags.
		if contains(r.TagCols, cols[i]) {
			// Tags are formatted as column name:value.
			tag := cols[i]
			switch s := c.(type) {
			case string: // varchar
				tags = append(tags, tag+":"+s)
			case []uint8: // char (fixed-length string)
				tags = append(tags, tag+":"+string(s))
			default:
				return metrics, errors.Errorf("failed to cast interface to string: col = %s, type = %T(%v)", cols[i], c, c)
			}
		}
	}

	// Build metrics by assigning tags for each value.
	for i, v := range values {
		metric := metric{
			name:  r.Name + "." + r.ValueCols[i],
			value: v,
			tags:  tags,
		}
		metrics = append(metrics, metric)
	}
	return metrics, nil
}

// contains returns true if the specified string is included in the list of strings.
func contains(s []string, key string) bool {
	for _, e := range s {
		if key == e {
			return true
		}
	}
	return false
}
