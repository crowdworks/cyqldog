package cyqldog

import (
	"fmt"

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

// NewStatsd returns an instance of statsd.Client.
func NewStatsd(d Dogstatsd) (*statsd.Client, error) {
	address := fmt.Sprintf("%s:%s", d.Host, d.Port)
	c, err := statsd.New(address)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open statsd: host=%s port=%s", d.Host, d.Port)
	}

	c.Namespace = d.Namespace + "."
	c.Tags = append(c.Tags, d.Tags...)
	return c, nil
}
