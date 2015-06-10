package cli_test

import (
	"github.com/qadium/plumber/cli"
	"io/ioutil"
	"reflect"
	"testing"
	"os"
)

// config with everything filled out
const goodConfig = `
language: python
name: foobar
inputs:
  - name: a
    description: none
    type: int
outputs:
  - name: b
    description: none
    type: int
env:
  - SECRET=1234
install:
  - make -j4
  - make install`

// config with all optional stuff left out
// includes an extra field which should just be ignored
const optConfig = `
language: python
name: foobar
inputs:
  - name: a
    type: string
outputs:
  - name: b
    type: string
extra: field`

// config with missing fields
const missingFields = `
language: go
name: foobar`

// config with missing name on inputs
const missingInputsAndOutputs = `
language: go
name: foobar
inputs:
  - type: foo
outputs:
  - foobar: a`

func writeConfig(t *testing.T, config string) string {
	configFile, err := ioutil.TempFile("", "plumberTest")
	if err != nil {
		t.Errorf("Had an issue creating the temp file to test, '%v'", err.Error())
	}

	if _, err := configFile.WriteString(config); err != nil {
		t.Errorf("Had an error writing the config file to test, '%v'", err.Error())
	}

	return configFile.Name()
}

// if the expectedContext is nil, then we'll check if an error is thrown
func parseConfig(t *testing.T, expectedContext *cli.Context, config string) {
	configFile := writeConfig(t, config)
	defer func() {
		if err := os.RemoveAll(configFile); err != nil {
			t.Errorf("Had an issue removing the temp file, '%v'", err.Error())
		}
	}()

	ctx, err := cli.ParseConfig(configFile)

	if expectedContext == nil {
		if err == nil {
			t.Errorf("Expected an error, but got '%v'", ctx)
		}
	} else {
		if err != nil {
			t.Errorf("Had an issue parsing the config file, '%v'", err.Error())
		}

		if !reflect.DeepEqual(ctx, expectedContext) {
			t.Errorf("Got '%v', expected '%v'", ctx, expectedContext)
		}
	}
}

func TestParseConfig(t *testing.T) {
	ctx := &cli.Context{
		Language: "python",
		Name: "foobar",
		Inputs: []cli.Field{
			cli.Field{
				Name: "a",
				Description: "none",
				Type: "int",
			},
		},
		Outputs: []cli.Field {
			cli.Field{
				Name: "b",
				Description: "none",
				Type: "int",
			},
		},
		Env: []string{"SECRET=1234"},
		Install: []string{"make -j4", "make install"},
	}
	parseConfig(t, ctx, goodConfig)
}

func TestParseOptConfig(t *testing.T) {
	ctx := &cli.Context {
		Language: "python",
		Name: "foobar",
		Inputs: []cli.Field{
			cli.Field{
				Name: "a",
				Description: "",
				Type: "string",
			},
		},
		Outputs: []cli.Field {
			cli.Field{
				Name: "b",
				Description: "",
				Type: "string",
			},
		},
		Env: nil,
		Install: nil,
	}
	parseConfig(t, ctx, optConfig)
}

func TestParseMissingFields(t *testing.T) {
	parseConfig(t, nil, missingFields)
}

func TestParseMissingInputsAndOutputs(t *testing.T) {
	parseConfig(t, nil, missingInputsAndOutputs)
}

func TestParseConfigFromDir(t *testing.T) {

}
