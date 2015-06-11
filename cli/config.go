package cli

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Field struct {
	Name        string
	Description string `yaml:",omitempty"`
	Type        string
}

type Bundle struct {
	Language string
	Name     string
	Inputs   []Field  `yaml:",flow"`
	Outputs  []Field  `yaml:",flow"`
	Env      []string `yaml:",flow,omitempty"`
	Install  []string `yaml:",flow,omitempty"`
}

const bundleConfig = ".plumb.yml"

// Parse a `.plumb.yml` in the given directory
func ParseBundleFromDir(path string) (*Bundle, error) {
	return ParseBundle(fmt.Sprintf("%s/%s", path, bundleConfig))
}

// Parse the config at the given path
func ParseBundle(path string) (*Bundle, error) {
	ctx := Bundle{}

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(bytes, &ctx); err != nil {
		return nil, err
	}

	if ctx.Language == "" {
		return nil, errors.New("You must provide a 'language' field.")
	}

	if ctx.Name == "" {
		return nil, errors.New("You must provide a 'name' field.")
	}

	if ctx.Inputs == nil {
		return nil, errors.New("You must provide 'inputs'.")
	}

	if ctx.Outputs == nil {
		return nil, errors.New("You must provide 'outputs'.")
	}

	// check inputs
	for _, input := range ctx.Inputs {
		if input.Name == "" {
			return nil, errors.New("You must provide a 'name' field for your inputs.")
		}
	}

	// check outputs
	for _, output := range ctx.Outputs {
		if output.Name == "" {
			return nil, errors.New("You must provide a 'name' field for your outputs.")
		}
	}

	return &ctx, nil
}
