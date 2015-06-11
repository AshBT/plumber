package cli_test

import (
	"fmt"
	"github.com/qadium/plumber/cli"
	"io/ioutil"
	"os"
	"os/user"
	"testing"
)

const testPlumberDir = "foo"

const testBootstrapDir = "boot"

const testKubeSuffix = "k8s"

// mock for cli context (used for testing)
// uses temp directories
func NewTestContext(t *testing.T) (*cli.Context, string) {
	// use the current user to store the plumb data
	usr, err := user.Current()
	if err != nil {
		t.Errorf("Got an error getting current user: '%v'", err)
	}

	tempDir, err := ioutil.TempDir(usr.HomeDir, "plumberTest")
	if err != nil {
		t.Errorf("Got an error constructing context: '%v'", err)
	}

	d := &cli.Context{
		fmt.Sprintf("%s/%s", tempDir, testPlumberDir),
		testKubeSuffix,
		"",
		"test-version",
		"plumber/test-manager",
		fmt.Sprintf("%s/%s", tempDir, testBootstrapDir),
	}
	return d, tempDir
}

func cleanTestDir(t *testing.T, tempDir string) {
	if err := os.RemoveAll(tempDir); err != nil {
		t.Errorf("Had an issue removing the temp file, '%v'", err)
	}
}

func TestPipelinePath(t *testing.T) {
	ctx, tempDir := NewTestContext(t)
	defer cleanTestDir(t, tempDir)

	expectedPath := fmt.Sprintf("%s/foobar", ctx.PipeDir)

	path := ctx.PipelinePath("foobar")
	if expectedPath != path {
		t.Error("PipelinePath: did not return expected path.")
	}
}

func TestGetPipeline(t *testing.T) {
	ctx, tempDir := NewTestContext(t)
	defer cleanTestDir(t, tempDir)

	expectedPath := fmt.Sprintf("%s/mypipe", ctx.PipeDir)
	// first, check that we fail with a "no such directory"
	path, err := ctx.GetPipeline("mypipe")
	if err == nil || err.Error() != fmt.Sprintf("stat %s: no such file or directory", expectedPath) {
		t.Errorf("We expected to fail with an error with no such directory, but got '%v' instead", err)
	}
	// make the expected directory and delete it after this test
	if err := os.MkdirAll(expectedPath, 0755); err != nil {
		t.Errorf("Encountered error making test directory '%s': '%v'", expectedPath, err)
	}
	defer os.RemoveAll(expectedPath)

	path, _ = ctx.GetPipeline("mypipe")
	if path != expectedPath {
		t.Error("GetPipeline: did not return expected path.")
	}
}

func TestKubernetesPath(t *testing.T) {
	ctx, tempDir := NewTestContext(t)
	defer cleanTestDir(t, tempDir)

	expectedPath := fmt.Sprintf("%s/barbaz/k8s", ctx.PipeDir)

	path := ctx.KubernetesPath("barbaz")
	if expectedPath != path {
		t.Error("KubernetesPath: did not return expected path.")
	}
}

func TestDefaultContext(t *testing.T) {
	usr, err := user.Current()
	if err != nil {
		t.Errorf("DefaultContext: Got an error getting current user: '%v'", err)
	}

	ctx, err := cli.NewDefaultContext()
	if err != nil {
		t.Errorf("DefaultContext: Got error '%v'", err)
	}

	if ctx.PipeDir != fmt.Sprintf("%s/.plumber", usr.HomeDir) ||
		ctx.KubeSuffix != "k8s" || ctx.Image != "plumber/manager" ||
		ctx.BootstrapDir != fmt.Sprintf("%s/.plumber-bootstrap", usr.HomeDir) {
		t.Errorf("DefaultContext: '%v' was not expected.", ctx)
	}

}
