package cli_test

import (
	"github.com/qadium/plumber/cli"
	"io/ioutil"
	"reflect"
	"testing"
	"os"
	"fmt"
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
const missingInputName = `
language: go
name: foobar
inputs:
  - type: foo
outputs:
  - foobar: a`

const missingOutputName = `
language: go
name: foobar
outputs:
  - hello: 1
inputs:
  - name: a`

func writeConfig(t *testing.T, config string) string {
	configFile, err := ioutil.TempFile("", "plumberTest")
	if err != nil {
		t.Errorf("Had an issue creating the temp file to test, '%v'", err)
	}

	if _, err := configFile.WriteString(config); err != nil {
		t.Errorf("Had an error writing the config file to test, '%v'", err)
	}

	return configFile.Name()
}

// if the expectedContext is nil, then we'll check if an error is thrown
func parseConfig(t *testing.T, expectedContext *cli.Context, config string) {
	configFile := writeConfig(t, config)
	defer func() {
		if err := os.RemoveAll(configFile); err != nil {
			t.Errorf("Had an issue removing the temp file, '%v'", err)
		}
	}()

	ctx, err := cli.ParseConfig(configFile)

	if expectedContext == nil {
		if err == nil {
			t.Errorf("Expected an error, but got '%v'", ctx)
		}
	} else {
		if err != nil {
			t.Errorf("Had an issue parsing the config file, '%v'", err)
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

func TestParseMissingInputName(t *testing.T) {
	parseConfig(t, nil, missingInputName)
}

func TestParseMissingOutputName(t *testing.T) {
	parseConfig(t, nil, missingOutputName)
}

func TestParseMissingLanguage(t *testing.T) {
	parseConfig(t, nil, `name: hello`)
}

func TestParseMissingOutputs(t *testing.T) {
	parseConfig(t, nil, `language: python
name: hello
inputs:
  - name: a`)
}
func TestParseMissingName(t *testing.T) {
	parseConfig(t, nil, `language: python`)
}

func TestParseConfigFromDir(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "plumberTest")
	if err != nil {
		t.Errorf("Could not make temp dir; got error '%v'", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Errorf("Had an issue removing the temp file, '%v'", err)
		}
	}()

	filename := fmt.Sprintf("%s/.plumb.yml", tempDir)

	if err := ioutil.WriteFile(filename, []byte(goodConfig), 0644); err != nil {
		t.Errorf("Could not write .plumb.yml; got error '%v'", err)
	}

	expected := &cli.Context{
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

	ctx, err := cli.ParseConfigFromDir(tempDir)
	if err != nil {
		t.Errorf("Should have parsed properly, got error '%v'", err)
	}

	if !reflect.DeepEqual(ctx, expected) {
		t.Errorf("Got '%v', expected '%v'", ctx, expected)
	}
}
