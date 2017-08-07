package cyqldog

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"os"
	"strings"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

// Config represents the structure of the configuration file.
type Config struct {
	// Db is a configuration of the database to connect.
	Db DataSource `yaml:"db"`
	// Statsd is a configuration of the dogstatsd to connect.
	Statd Dogstatsd `yaml:"dogstatsd"`
	// Rules are a list of rules to monitor
	Rules []Rule `yaml:"rules"`
}

// NewConfig returns an instance of the Config.
func NewConfig(filename string) (*Config, error) {
	// Read bytes from file.
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read file: %s", filename)
	}

	// Render envionment variables in the configuration file.
	rendered, err := renderEnv(buf)
	if err != nil {
		return nil, err
	}

	// Parse yaml into Config.
	c := Config{}
	err = yaml.Unmarshal(rendered, &c)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse yaml: %s", filename)
	}

	return &c, nil
}

// envMap is a map of environment variables.
type envMap map[string]string

// newEnvMap returns a map of environment variables.
func newEnvMap() *envMap {
	envs := make(envMap)
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		envs[pair[0]] = pair[1]
	}
	return &envs
}

// renderEnv embeds environment variables with input buffer as a template.
func renderEnv(buf []byte) ([]byte, error) {
	tmpl, err := template.New("env").Parse(string(buf))
	if err != nil {
		return []byte{}, errors.Wrapf(err, "failed to parse template: %v", buf)
	}

	var rendered bytes.Buffer
	envs := newEnvMap()
	err = tmpl.Execute(&rendered, *envs)
	if err != nil {
		return []byte{}, errors.Wrapf(err, "failed to execute template: %v", buf)
	}
	return rendered.Bytes(), nil
}
