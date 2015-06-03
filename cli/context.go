package cli

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Field struct {
	Name        string
	Description string `yaml:",omitempty"`
	Type        string
}

type Context struct {
	Language string
	Name     string
	Inputs   []Field  `yaml:",flow"`
	Outputs  []Field  `yaml:",flow"`
	Env      []string `yaml:",flow,omitempty"`
	Install  []string `yaml:",flow,omitempty"`
}

const bundleConfig = ".plumb.yml"

func parseConfigFromDir(path string) (*Context, error) {
	return parseConfig(fmt.Sprintf("%s/%s", path, bundleConfig))
}

func parseConfig(path string) (*Context, error) {
	ctx := Context{}

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(bytes, &ctx); err != nil {
		return nil, err
	}
	return &ctx, nil
}
