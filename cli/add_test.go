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
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestAddExistingPipeline(t *testing.T) {
	ctx, tempDir := NewTestContext(t)
	defer cleanTestDir(t, tempDir)

	// create a pipeline
	if err := ctx.Create("add-test"); err != nil {
		t.Errorf("TestAdd: error creating '%v'", err)
	}
	project := fmt.Sprintf("%s/add-test", ctx.PipeDir)

	// make a usable bundle
	createTestBundle(t, "barbaz", tempDir)

	// add that bundle to the pipeline
	if err := ctx.Add("add-test", tempDir); err != nil {
		t.Errorf("TestAdd: '%v'", err)
	}

	// check that the directory now contains a 'barbaz.yml'
	if _, err := os.Stat(fmt.Sprintf("%s/barbaz.yml", project)); err != nil {
		t.Errorf("TestAdd: config file did not get copied over, '%v'", err)
	}

	// check that git history has advanced
	cmd := exec.Command("git", "-C", project, "log", "-1", "--pretty=%s")
	bytes, err := cmd.Output()
	if err != nil {
		t.Errorf("TestAdd: got an error: '%v'", err)
	}
	output := strings.Trim(string(bytes), "\n")
	if output != "Updated 'barbaz' config." {
		t.Errorf("TestAdd: expected git config to be updated, but got '%s'", output)
	}
}

func TestAddWithoutPipeline(t *testing.T) {
	ctx, tempDir := NewTestContext(t)
	defer cleanTestDir(t, tempDir)

	// add nonexistent bundle to nonexistent pipeline
	err := ctx.Add("add-test-nonexistent", tempDir)
	project := fmt.Sprintf("%s/add-test-nonexistent", ctx.PipeDir)

	if err == nil || err.Error() != fmt.Sprintf("stat %s: no such file or directory", project) {
		t.Errorf("TestAddWithoutPipeline: did not get an expected error!")
	}
}

func TestAddWithoutBundle(t *testing.T) {
	ctx, tempDir := NewTestContext(t)
	defer cleanTestDir(t, tempDir)

	// create a pipeline
	if err := ctx.Create("add-test-no-bundle"); err != nil {
		t.Errorf("TestAddWithoutBundle: '%v'", err)
	}

	// add nonexitent bundle to the pipeline (should give an error)
	err := ctx.Add("add-test-no-bundle", tempDir)
	if err == nil || err.Error() != fmt.Sprintf("open %s/.plumb.yml: no such file or directory", tempDir) {
		t.Errorf("TestAddWithoutBundle: did not get an expected error!")
	}
}
