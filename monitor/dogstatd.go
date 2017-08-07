package monitor

import (
	"fmt"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/pkg/errors"
)

type Dogstatsd struct {
	Host      string   `yaml:"host"`
	Port      string   `yaml:"port"`
	Namespace string   `yaml:"namespace"`
	Tags      []string `yaml:"tags"`
}

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
