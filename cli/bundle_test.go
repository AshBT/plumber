/**
 * Copyright 2015 Qadium, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package cli_test

import (
	"bytes"
	"fmt"
	"github.com/qadium/plumber/shell"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

// TODO: this test code should live in the language-dependent repository
// and our test should *load* those tests
//
// this test echoes the input
const bundleTestPython = `
def run(in_dict):
	in_dict['b'] = "echo {}".format(in_dict['a'])
	return in_dict
`

const bundleFormat = `
language: python
name: %s
inputs:
  - name: a
    type: string
outputs:
  - name: b
    type: string
extra: field`

func createTestBundle(t *testing.T, bundleName, tempDir string) {
	filename := fmt.Sprintf("%s/.plumb.yml", tempDir)
	bundle := fmt.Sprintf(bundleFormat, bundleName)
	if err := ioutil.WriteFile(filename, []byte(bundle), 0644); err != nil {
		t.Errorf("Could not write .plumb.yml; got error '%v'", err)
	}

	// write out the python test code
	filename = fmt.Sprintf("%s/%s.py", tempDir, bundleName)
	if err := ioutil.WriteFile(filename, []byte(bundleTestPython), 0644); err != nil {
		t.Errorf("Could not write foobar.py; got error '%v'", err)
	}

	// write out the empty requirements.txt file
	filename = fmt.Sprintf("%s/requirements.txt", tempDir)
	if err := ioutil.WriteFile(filename, []byte{}, 0644); err != nil {
		t.Errorf("Could not write requirements.txt; got error '%v'", err)
	}
}

func TestBundle(t *testing.T) {
	//log.SetOutput(ioutil.Discard)

	ctx, tempDir := NewTestContext(t)
	defer cleanTestDir(t, tempDir)

	// create a bundle for testing
	createTestBundle(t, "foobar", tempDir)

	// bundle the tempDir
	err := ctx.Bundle(tempDir)
	if err != nil {
		t.Errorf("Bundle: Got an unxpected error while bundling, '%v'", err)
	}
	defer shell.RunAndLog(ctx.DockerCmd, "rmi", ctx.GetImage("foobar"))

	// run the container and check that it increments the input
	if err := shell.RunAndLog(ctx.DockerCmd, "run", "-d", "-p", "9800:9800", "--name", "foobar", ctx.GetImage("foobar")); err != nil {
		t.Errorf("Bundle: Got an error during docker run: '%v'", err)
	}
	defer shell.RunAndLog(ctx.DockerCmd, "rm", "-f", "foobar")
	// wait a bit for the container to come up
	time.Sleep(1 * time.Second)

	// step 4. send some JSON and check for echos
	// first, find the IP to connect to
	hostIp, err := ctx.GetDockerHost()
	if err != nil {
		t.Errorf("Bootstrap: Got an error getting the docker host: '%v'", err)
	}

	// second, send over some JSON and verify result
	resp, err := http.Post(fmt.Sprintf("http://%s:9800", hostIp), "application/json", bytes.NewBufferString(`{"a": "trusty"}`))
	if err != nil {
		t.Error(err)
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	result := buf.String()
	if result != `{"a": "trusty", "b": "echo trusty"}` {
		t.Errorf("Bundle: Got '%s'; did not get expected response", result)
	}
}
