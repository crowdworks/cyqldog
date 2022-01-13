package cyqldog

import (
	"bytes"
	"golang.org/x/xerrors"
	"html/template"
	"io/ioutil"
	"os"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

// Config represents the structure of the configuration file.
type Config struct {
	// DB is a configuration of the database to connect.
	DB DataSourceConfig `yaml:"data_source"`
	// Notifiers are configurations of output plugins.
	Notifiers NotifiersConfig `yaml:"notifiers"`
	// Rules are a list of rules to monitor
	Rules []Rule `yaml:"rules"`
}

// newConfig returns an instance of the Config.
func newConfig(filename string) (*Config, error) {
	// Read bytes from file.
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, xerrors.Errorf("failed to read file: %s: %w", filename, err)
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
		return nil, xerrors.Errorf("failed to parse yaml: %s: %w", filename, err)
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
		return []byte{}, xerrors.Errorf("failed to parse template: %v: %w", buf, err)
	}

	var rendered bytes.Buffer
	envs := newEnvMap()
	err = tmpl.Execute(&rendered, *envs)
	if err != nil {
		return []byte{}, xerrors.Errorf("failed to execute template: %v: %w", buf, err)
	}
	return rendered.Bytes(), nil
}
