package cli_test

import (
	"fmt"
	"github.com/qadium/plumber/cli"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

// bundle with everything filled out
const goodBundle = `
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

// bundle with all optional stuff left out
// includes an extra field which should just be ignored
const optBundle = `
language: python
name: foobar
inputs:
  - name: a
    type: string
outputs:
  - name: b
    type: string
extra: field`

// bundle with missing fields
const missingFields = `
language: go
name: foobar`

// bundle with missing name on inputs
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

func writeBundle(t *testing.T, bundle string) string {
	configFile, err := ioutil.TempFile("", "plumberTest")
	if err != nil {
		t.Errorf("Had an issue creating the temp file to test, '%v'", err)
	}

	if _, err := configFile.WriteString(bundle); err != nil {
		t.Errorf("Had an error writing the bundle file to test, '%v'", err)
	}

	return configFile.Name()
}

// if the expectedBundle is nil, then we'll check if an error is thrown
func parseBundle(t *testing.T, expectedBundle *cli.Bundle, bundle string) {
	configFile := writeBundle(t, bundle)
	defer func() {
		if err := os.RemoveAll(configFile); err != nil {
			t.Errorf("Had an issue removing the temp file, '%v'", err)
		}
	}()

	ctx, err := cli.ParseBundle(configFile)

	if expectedBundle == nil {
		if err == nil {
			t.Errorf("Expected an error, but got '%v'", ctx)
		}
	} else {
		if err != nil {
			t.Errorf("Had an issue parsing the bundle file, '%v'", err)
		}

		if !reflect.DeepEqual(ctx, expectedBundle) {
			t.Errorf("Got '%v', expected '%v'", ctx, expectedBundle)
		}
	}
}

func TestParseBundle(t *testing.T) {
	ctx := &cli.Bundle{
		Language: "python",
		Name:     "foobar",
		Inputs: []cli.Field{
			cli.Field{
				Name:        "a",
				Description: "none",
				Type:        "int",
			},
		},
		Outputs: []cli.Field{
			cli.Field{
				Name:        "b",
				Description: "none",
				Type:        "int",
			},
		},
		Env:     []string{"SECRET=1234"},
		Install: []string{"make -j4", "make install"},
	}
	parseBundle(t, ctx, goodBundle)
}

func TestParseOptBundle(t *testing.T) {
	ctx := &cli.Bundle{
		Language: "python",
		Name:     "foobar",
		Inputs: []cli.Field{
			cli.Field{
				Name:        "a",
				Description: "",
				Type:        "string",
			},
		},
		Outputs: []cli.Field{
			cli.Field{
				Name:        "b",
				Description: "",
				Type:        "string",
			},
		},
		Env:     nil,
		Install: nil,
	}
	parseBundle(t, ctx, optBundle)
}

func TestParseMissingFields(t *testing.T) {
	parseBundle(t, nil, missingFields)
}

func TestParseMissingInputName(t *testing.T) {
	parseBundle(t, nil, missingInputName)
}

func TestParseMissingOutputName(t *testing.T) {
	parseBundle(t, nil, missingOutputName)
}

func TestParseMissingLanguage(t *testing.T) {
	parseBundle(t, nil, `name: hello`)
}

func TestParseMissingOutputs(t *testing.T) {
	parseBundle(t, nil, `language: python
name: hello
inputs:
  - name: a`)
}
func TestParseMissingName(t *testing.T) {
	parseBundle(t, nil, `language: python`)
}

func TestParseBundleFromDir(t *testing.T) {
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

	if err := ioutil.WriteFile(filename, []byte(goodBundle), 0644); err != nil {
		t.Errorf("Could not write .plumb.yml; got error '%v'", err)
	}

	expected := &cli.Bundle{
		Language: "python",
		Name:     "foobar",
		Inputs: []cli.Field{
			cli.Field{
				Name:        "a",
				Description: "none",
				Type:        "int",
			},
		},
		Outputs: []cli.Field{
			cli.Field{
				Name:        "b",
				Description: "none",
				Type:        "int",
			},
		},
		Env:     []string{"SECRET=1234"},
		Install: []string{"make -j4", "make install"},
	}

	ctx, err := cli.ParseBundleFromDir(tempDir)
	if err != nil {
		t.Errorf("Should have parsed properly, got error '%v'", err)
	}

	if !reflect.DeepEqual(ctx, expected) {
		t.Errorf("Got '%v', expected '%v'", ctx, expected)
	}
}
