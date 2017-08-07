package monitor

import (
	"time"

	"github.com/pkg/errors"
)

type Rule struct {
	Name      string        `yaml:"name"`
	Interval  time.Duration `yaml:"interval"`
	Query     string        `yaml:"query"`
	ValueCols []string      `yaml:"value_cols"`
	TagCols   []string      `yaml:"tag_cols"`
}

func (r *Rule) buildMetrics(row []interface{}, cols []string) ([]metric, error) {
	metrics := []metric{}
	values := []float64{}
	tags := []string{}

	for i, c := range row {
		if contains(r.ValueCols, cols[i]) {
			switch f := c.(type) {
			case float64:
				values = append(values, f)
			case int64:
				values = append(values, float64(f))
			default:
				return metrics, errors.Errorf("failed to cast from interface to float64: col = %s, type = %T(%v)", cols[i], c, c)
			}
		}

		if contains(r.TagCols, cols[i]) {
			tag := cols[i]
			switch s := c.(type) {
			case string:
				tags = append(tags, tag+":"+s)
			case []uint8:
				tags = append(tags, tag+":"+string(s))
			default:
				return metrics, errors.Errorf("failed to cast interface to string: col = %s, type = %T(%v)", cols[i], c, c)
			}
		}
	}

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

func contains(s []string, key string) bool {
	for _, e := range s {
		if key == e {
			return true
		}
	}
	return false
}
