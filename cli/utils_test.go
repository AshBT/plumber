package cli_test

import (
	"fmt"
	"github.com/qadium/plumber/cli"
	"os"
	"os/user"
	"testing"
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
	if err == nil || err.Error() != fmt.Sprintf("stat %s: no such file or directory", expectedPath) {
		t.Errorf("We expected to fail with an error with no such directory, but got '%v' instead", err)
	}
	// make the expected directory and delete it after this test
	if err := os.MkdirAll(expectedPath, 0755); err != nil {
		t.Errorf("Encountered error making test directory '%s': '%v'", expectedPath, err)
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
