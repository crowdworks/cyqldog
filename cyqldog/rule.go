package cyqldog

import "time"

// Rule represents monitoring conditions.
type Rule struct {
	// Name of the rule.
	// The metic name sent to the dogstatsd is:
	//   Dogstatsd.Namespace + Rule.Name + Rule.ValueCols[*]
	Name string `yaml:"name"`
	// Interval of the monitoring.
	Interval time.Duration `yaml:"interval"`
	// Query to the database.
	Query string `yaml:"query"`
	// Notifier is a name of notifier to send metrics.
	Notifier string `yaml:"notifier"`
	// ValueCols is a list of names of the columns used as metric values.
	ValueCols []string `yaml:"value_cols"`
	// TagCols is a list of names of the columns used as metric tags.
	TagCols []string `yaml:"tag_cols"`
}
