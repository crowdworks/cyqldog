package monitor

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"os"
	"strings"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Db    DataSource `yaml:"db"`
	Statd Dogstatsd  `yaml:"dogstatsd"`
	Rules []Rule     `yaml:"rules"`
}

func NewConfig(filename string) (*Config, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read file: %s", filename)
	}

	rendered, err := renderEnv(buf)
	if err != nil {
		return nil, err
	}

	c := Config{}
	err = yaml.Unmarshal(rendered, &c)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse yaml: %s", filename)
	}

	return &c, nil
}

type envMap map[string]string

func newEnvMap() *envMap {
	envs := make(envMap)
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		envs[pair[0]] = pair[1]
	}
	return &envs
}

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
