package cyqldog

// Notifier is an interface which send metrics to.
type Notifier interface {
	Put(qr QueryResult, rule Rule) error
}

// Notifiers is a map of notifiers.
type Notifiers map[string]Notifier

// NotifiersConfig are configurations of notifiers.
type NotifiersConfig struct {
	// Dogstatsd is a configuration of the dogstatsd to connect.
	Dogstatsd DogstatsdConfig `yaml:"dogstatsd"`
}

// newNotifiers returns an instance of Notifiers.
func newNotifiers(c NotifiersConfig) (Notifiers, error) {
	notifiers := make(Notifiers)

	dogstatsd, err := newDogstatsd(c.Dogstatsd)
	if err != nil {
		return notifiers, err
	}

	notifiers["dogstatsd"] = dogstatsd

	return notifiers, nil
}
