package cli_test

import (
	"bytes"
	"fmt"
	"github.com/qadium/plumber/cli"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"os/user"
	"strings"
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
		"manager",
		fmt.Sprintf("%s/%s", tempDir, testBootstrapDir),
		"plumber_test",
	}
	return d, tempDir
}

func cleanTestDir(t *testing.T, tempDir string) {
	if err := os.RemoveAll(tempDir); err != nil {
		t.Errorf("Had an issue removing the temp directory, '%v'", err)
	}
}

func getImageIp(t *testing.T, imageName string) string {
	// get the DOCKER_HOST environment variable; if not defined, use
	// docker to find it
	hostIp := os.Getenv("DOCKER_HOST")
	if hostIp == "" {
		cmd := exec.Command("docker", "inspect", "--format='{{.NetworkSettings.Gateway}}'", imageName)
		hostIpBytes, err := cmd.Output()
		if err != nil {
			t.Errorf("Got an error during docker inspect: '%v'", err)
		}
		hostIpBytes = bytes.Trim(hostIpBytes, "\r\n")
		hostIp = string(hostIpBytes)
	} else {
		hostUrl, err := url.Parse(hostIp)
		if err != nil {
			t.Errorf("Got an error during url parsing: '%v'", err)
		}
		// docker host is usually in the form of IP:PORT
		hostIp = strings.Split(hostUrl.Host, ":")[0]
	}
	return hostIp
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
		ctx.KubeSuffix != "k8s" || ctx.ManagerImage != "manager" ||
		ctx.BootstrapDir != fmt.Sprintf("%s/.plumber-bootstrap", usr.HomeDir) ||
		ctx.ImageRepo != "plumber" {
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
