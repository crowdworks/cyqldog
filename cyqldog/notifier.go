package cyqldog

import "fmt"

// Notifier is an interface which send metrics to.
type Notifier interface {
	Put(qr QueryResult, rule Rule) error
	Event(e *Event) error
}

// An Event is an object that can be posted to the Notifier.
type Event struct {
	// Title of the event. Required.
	Title string
	// Text is the description of the event. Required.
	Text string
	// Level is a level of the event.
	// This can be info, error, warning or success. (default: info)
	Level string
	// Tags for the event.
	Tags []string
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

func newErrorEvent(err error) *Event {
	return &Event{
		Title: fmt.Sprintf("cyqldog: %s", err),
		Text:  fmt.Sprintf("%+v", err),
		Level: "error",
		Tags:  []string{"cyqldog"},
	}
}
