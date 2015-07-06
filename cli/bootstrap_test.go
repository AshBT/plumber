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
	"net/http"
	"testing"
	"time"
	"os"
)

// this is a *functional test*
// we check that boostrap works by actually running the boostrap command
// and checking that the container is built and runs
func TestBootstrap(t *testing.T) {
	ctx, tempDir := NewTestContext(t)
	defer cleanTestDir(t, tempDir)

	// step 1. remove any image named plumber/test-manager from the current
	// set of docker images (ignore any errors)
	_ = shell.RunAndLog(ctx.DockerCmd, "rmi", ctx.GetManagerImage())

	// step 2. invoke Bootstrap for building ctx.GetManagerImage()
	if err := ctx.Bootstrap(); err != nil {
		t.Errorf("Bootstrap: Got an error during bootstrap: '%v'", err)
	}
	// remove the manager image if we're on travis
	// note that this is a "hack" to avoid space requirements on travis
	// since we can bootstrap once, but not twice?
	if os.Getenv("TRAVIS") == "" {
		defer shell.RunAndLog(ctx.DockerCmd, "rmi", ctx.GetManagerImage())
	}

	// if on travis, try to remove the golang container
	if os.Getenv("TRAVIS") != "" {
		if err := shell.RunAndLog(ctx.DockerCmd, "rmi", "centurylink/golang-builder"); err != nil {
			t.Errorf("Bootstrap: Got an error removing golang-builder: '%v'", err)
		}
	}

	// step 3. run the image (it *should* just echo in response)
	if err := shell.RunAndLog(ctx.DockerCmd, "run", "-d", "-p", "9800:9800", "--name", "plumber-test", ctx.GetManagerImage()); err != nil {
		t.Errorf("Bootstrap: Got an error during docker run: '%v'", err)
	}
	defer shell.RunAndLog(ctx.DockerCmd, "rm", "-f", "plumber-test")
	// wait a bit for the container to come up
	time.Sleep(1 * time.Second)

	// step 4. send some JSON and check for echos
	hostIp, err := ctx.GetDockerHost()
	if err != nil {
		t.Errorf("Bootstrap: Got an error getting the docker host: '%v'", err)
	}

	// second, send over some JSON and verify result
	resp, err := http.Post(fmt.Sprintf("http://%s:9800", hostIp), "application/json", bytes.NewBufferString(`{"foo": 3}`))
	if err != nil {
		t.Error(err)
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	result := buf.String()
	if result != `{"foo": 3}` {
		t.Errorf("Bootstrap: Got '%s'; did not get expected response", result)
	}
}
