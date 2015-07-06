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
	"github.com/qadium/plumber/cli"
	"github.com/qadium/plumber/shell"
	"github.com/qadium/plumber/graph"
	"net/http"
	"strings"
	"syscall"
	"testing"
	"time"
	"os"
)

func createTestBundleAndPipeline(t *testing.T, ctx *cli.Context, pipeline, bundleName, tempDir string) {
	// bootstrap the manager if not on travis
	if os.Getenv("TRAVIS") == "" {
		if err := ctx.Bootstrap(); err != nil {
			t.Errorf("CreateTestBundleAndPipeline: Got an error during bootstrap: '%v'", err)
		}
	}

	// create a pipeline
	if err := ctx.Create(pipeline); err != nil {
		t.Errorf("CreateTestBundleAndPipeline: error creating '%v'", err)
	}

	// make a usable bundle and bundle it
	createTestBundle(t, bundleName, tempDir)
	if err := ctx.Bundle(tempDir); err != nil {
		t.Errorf("CreateTestBundleAndPipeline: error bundling test bundle, '%v'", err)
	}

	// add that bundle to the pipeline
	if err := ctx.Add(pipeline, tempDir); err != nil {
		t.Errorf("CreateTestBundleAndPipeline: '%v'", err)
	}
}

// Tests the Start command.
func TestStart(t *testing.T) {
	ctx, tempDir := NewTestContext(t)
	defer cleanTestDir(t, tempDir)

	const testPipeline = "test-start"
	const testBundle = "bazbux"
	createTestBundleAndPipeline(t, ctx, testPipeline, testBundle, tempDir)
	defer shell.RunAndLog("docker", "rmi", ctx.GetImage(testBundle))
	defer shell.RunAndLog("docker", "rmi", ctx.GetManagerImage())

	// set the interrupt handler to go off after 5 seconds
	go func() {
		time.Sleep(5 * time.Second)
		syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	}()

	// send a post request to the "server" and see what we get back
	go func() {
		time.Sleep(4 * time.Second)
		hostIp, err := ctx.GetDockerHost()
		if err != nil {
			t.Errorf("TestStart: Got an error getting the docker host: '%v'", err)
		}
		resp, err := http.Post(fmt.Sprintf("http://%s:9800", hostIp), "application/json", bytes.NewBufferString(`{"data": {"a": "trusty"}}`))
		if err != nil {
			t.Error(err)
		} else {
			buf := new(bytes.Buffer)
			buf.ReadFrom(resp.Body)
			result := buf.String()
			if !strings.Contains(result, `{"a": "trusty", "b": "echo trusty"}`) {
				t.Errorf("TestStart: Got '%s'; did not contain expected response", result)
			}
		}
	}()

	// start the pipeline locally (set the gce project to '' to run
	// locally)
	err := ctx.Start(testPipeline, "")
	if err != nil {
		t.Errorf("TestStart: '%v'", err)
	}

	// now attempt to start it remotely
	const projectId = "gce-project-id"
	err = ctx.Start(testPipeline, projectId)
	if err != nil {
		t.Errorf("TestStart: [remote] '%v'", err)
	}
	remoteImage := fmt.Sprintf("gcr.io/%s/plumber-%s", projectId, "manager")
	defer shell.RunAndLog("docker", "rmi", remoteImage)
	remoteImage = fmt.Sprintf("gcr.io/%s/plumber-%s", projectId, testBundle)
	defer shell.RunAndLog("docker", "rmi", remoteImage)
}

func TestStartNonExistentPipeline(t *testing.T) {
	ctx, tempDir := NewTestContext(t)
	defer cleanTestDir(t, tempDir)

	err := ctx.Start("", "")
	if err == nil || err.Error() != fmt.Sprintf("stat %s/: no such file or directory", ctx.PipeDir) {
		t.Errorf("TestStartNonExistentPipeline: did not get expected error, '%v'", err)
	}
}

func TestBundlesToGraph(t *testing.T) {
	// type Bundle struct {
	// 	Language string
	// 	Name     string
	// 	Inputs   []Field  `yaml:",flow"`
	// 	Outputs  []Field  `yaml:",flow"`
	// 	Env      []string `yaml:",flow,omitempty"`
	// 	Install  []string `yaml:",flow,omitempty"`
	// }
	bundles := []*cli.Bundle{
		&cli.Bundle{"python", "a", "",
			[]cli.Field{cli.Field{"text", "", "",},},
			[]cli.Field{cli.Field{"dude", "", "",},},
			nil, nil, nil,
		},
		&cli.Bundle{"python", "b", "",
			[]cli.Field{cli.Field{"id", "", "",},},
			[]cli.Field{cli.Field{"text", "", "",},},
			nil, nil, nil,
		},
		&cli.Bundle{"python", "c", "",
			[]cli.Field{cli.Field{"text", "", "",},},
			[]cli.Field{cli.Field{"bar", "", "",},},
			nil, nil, nil,
		},
		&cli.Bundle{"python", "d", "",
			[]cli.Field{cli.Field{"id", "", "",},cli.Field{"text", "", "",}, },
			[]cli.Field{cli.Field{"foo", "", "",},},
			nil, nil, nil,
		},
	}

	nodes := cli.BundlesToGraph(bundles)
	sorted, err := graph.ReverseTopoSort(nodes)
	if err != nil {
		t.Errorf("BundlesToGraph: Unexpected error '%v'", err)
	}
	if sorted[0] != "a" || sorted[1] != "c" || sorted[2] != "d" || sorted[3] != "b" {
		t.Errorf("BundlesToGraph: Did not properly handle dependencies, reverse sorted: '%v'.", sorted)
	}

}
