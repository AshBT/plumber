package cli_test

import (
	"fmt"
	"github.com/qadium/plumber/cli"
	"io/ioutil"
	"os"
	"os/user"
	"testing"
	"runtime"
)

const testPlumberDir = "foo"

const testBootstrapDir = "boot"

const testKubeSubdir = "k8s"

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
		testKubeSubdir,
		"",
		"test-version",
		"manager",
		fmt.Sprintf("%s/%s", tempDir, testBootstrapDir),
		"plumber_test",
		"docker",
		"docker0",
		"DOCKER_HOST",
		"true",
		"true",
	}
	return d, tempDir
}

func cleanTestDir(t *testing.T, tempDir string) {
	if err := os.RemoveAll(tempDir); err != nil {
		t.Errorf("Had an issue removing the temp directory, '%v'", err)
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
		ctx.KubeSubdir != "k8s" || ctx.ManagerImage != "manager" ||
		ctx.BootstrapDir != fmt.Sprintf("%s/.plumber-bootstrap", usr.HomeDir) ||
		ctx.ImageRepo != "plumber" || ctx.DockerCmd != "docker" ||
		ctx.DockerIface != "docker0" || ctx.DockerHostEnv != "DOCKER_HOST" ||
		ctx.GcloudCmd != "gcloud" || ctx.KubectlCmd != "kubectl" {
		t.Errorf("DefaultContext: '%v' was not expected.", ctx)
	}
}

func TestGetManagerImage(t *testing.T) {
	ctx, tempDir := NewTestContext(t)
	defer cleanTestDir(t, tempDir)

	imageName := ctx.GetManagerImage()
	if imageName != "plumber_test/manager" {
		t.Error("GetManagerImage: did not return expected image name.")
	}
}

func TestGetImage(t *testing.T) {
	ctx, tempDir := NewTestContext(t)
	defer cleanTestDir(t, tempDir)

	imageName := ctx.GetImage("whatnot")
	if imageName != "plumber_test/whatnot" {
		t.Error("GetImage: did not return expected image name.")
	}
}

func TestGetImageWithNoImageRepo(t *testing.T) {
	ctx, tempDir := NewTestContext(t)
	ctx.ImageRepo = ""
	defer cleanTestDir(t, tempDir)

	imageName := ctx.GetImage("IAmReallyHere")
	if imageName != "IAmReallyHere" {
		t.Error("GetImageWithNoImageRepo: did not return expected image name.")
	}
}

func TestGetDockerHostFail(t *testing.T) {
	ctx, tempDir := NewTestContext(t)
	defer cleanTestDir(t, tempDir)
	ctx.DockerHostEnv = "I_AM_AN_ENV_THAT_DOESNT_EXIST_***"
	ctx.DockerIface = "reallyYouHaveAnIfaceWithThisName?"

	hostIp, err := ctx.GetDockerHost()
	if hostIp != "" || err == nil {
		t.Errorf("GetDockerHostFail: did not fail as expected")
	}
	if err.Error() != "no such network interface" {
		t.Errorf("GetDockerHostFail: got an unexpected error '%v'", err)
	}
}

func TestGetDockerHostWithDockerHostEnv(t *testing.T) {
	ctx, tempDir := NewTestContext(t)
	defer cleanTestDir(t, tempDir)
	if err := os.Setenv("PLUMBER_TEST_ENV", "http://127.0.0.1"); err != nil {
		t.Errorf("GetDockerHostWithDockerHostEnv: did not set test env variable, '%v'.", err)
	}
	defer os.Unsetenv("PLUMBER_TEST_ENV")
	ctx.DockerHostEnv = "PLUMBER_TEST_ENV"

	hostIp, err := ctx.GetDockerHost()
	if err != nil {
		t.Errorf("GetDockerHostWithDockerHostEnv: got unexpected error '%v'.", err)
	}
	if hostIp != "127.0.0.1" {
		t.Errorf("GetDockerHostWithDockerHostEnv: did not get expected IP. Got '%s' instead.", hostIp)
	}
}

func TestGetDockerHostWithDockerIface(t *testing.T) {
	ctx, tempDir := NewTestContext(t)
	defer cleanTestDir(t, tempDir)
	ctx.DockerHostEnv = "I_AM_AN_ENV_THAT_DOESNT_EXIST_***"
	if runtime.GOOS == "darwin" {
		ctx.DockerIface = "lo0"
	} else if runtime.GOOS == "linux" {
		ctx.DockerIface = "lo"
	} else {
		t.Skipf("GetDockerHostWithDockerIface: skipping test for this os '%s'", runtime.GOOS)
	}

	hostIp, err := ctx.GetDockerHost()
	if err != nil {
		t.Errorf("GetDockerHostWithDockerIface: got unexpected error '%v'.", err)
	}
	if hostIp != "127.0.0.1" {
		t.Errorf("GetDockerHostWithDockerIface: did not get expected IP.")
	}
}
