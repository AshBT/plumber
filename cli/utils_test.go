package cli_test

import (
	"testing"
	"os"
	"os/user"
	"fmt"
	"github.com/qadium/plumber/cli"
)

func TestPipelinePath(t *testing.T) {
	usr, err := user.Current()
	if err != nil {
		t.Error("Could not get the user running this test. /shrug")
	}

	expectedPath := fmt.Sprintf("%s/.plumber/foobar", usr.HomeDir)

	path, err := cli.PipelinePath("foobar")
	if expectedPath != path {
		t.Error("PipelinePath: did not return expected path.")
	}
}

func TestGetPipeline(t *testing.T) {
	usr, err := user.Current()
	if err != nil {
		t.Error("Could not get the user running this test. /shrug")
	}

	expectedPath := fmt.Sprintf("%s/.plumber/you-better-hope-nobody-names-pipelines-this-way", usr.HomeDir)
	// first, check that we fail with a "no such directory"
	path, err := cli.GetPipeline("you-better-hope-nobody-names-pipelines-this-way")
	if err == nil || err.Error() != "stat /Users/echu/.plumber/you-better-hope-nobody-names-pipelines-this-way: no such file or directory" {
		t.Error("We expected to fail with an error with no such directory, but got '%v' instead", err.Error())
	}
	// make the expected directory and delete it after this test
	if err := os.MkdirAll(expectedPath, 0755); err != nil {
		t.Errorf("Encountered error making test directory '%s': '%v'", expectedPath, err.Error())
	}
	defer os.RemoveAll(expectedPath)

	path, err = cli.GetPipeline("you-better-hope-nobody-names-pipelines-this-way")
	if path != expectedPath {
		t.Error("GetPipeline: did not return expected path.")
	}
}

func TestKubernetesPath(t *testing.T) {
	usr, err := user.Current()
	if err != nil {
		t.Error("Could not get the user running this test. /shrug")
	}

	expectedPath := fmt.Sprintf("%s/.plumber/barbaz/k8s", usr.HomeDir)

	path, err := cli.KubernetesPath("barbaz")
	if expectedPath != path {
		t.Error("KubernetesPath: did not return expected path.")
	}
}
