package plumb

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"fmt"
)

type Field struct {
	Name        string
	Description string `yaml:",omitempty"`
	Type        string
}

type PlumbContext struct {
	Language string
	Name     string
	Inputs   []Field  `yaml:",flow"`
	Outputs  []Field  `yaml:",flow"`
	Env      []string `yaml:",flow,omitempty"`
	Install  []string `yaml:",flow,omitempty"`
}

const bundleConfig = ".plumb.yml"

func parseConfig(path string) (* PlumbContext, error) {
	ctx := PlumbContext{}

	bytes, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", path, bundleConfig))
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(bytes, &ctx); err != nil {
		return nil, err
	}
	return &ctx, nil
}
