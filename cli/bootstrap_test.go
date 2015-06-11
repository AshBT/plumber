package cli_test

import (
	"bytes"
	"fmt"
	"github.com/qadium/plumber/shell"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

// this is a *functional test*
// we check that boostrap works by actually running the boostrap command
// and checking that the container is built and runs

func TestBootstrap(t *testing.T) {
	ctx, tempDir := NewTestContext(t)
	defer cleanTestDir(t, tempDir)

	// step 1. remove any image named plumber/test-manager from the current
	// set of docker images (ignore any errors)
	_ = shell.RunAndLog("docker", "rmi", ctx.Image)

	// step 2. invoke Bootstrap for building ctx.Image
	if err := ctx.Bootstrap(); err != nil {
		t.Errorf("Got an error during bootstrap: '%v'", err)
	}
	defer shell.RunAndLog("docker", "rmi", ctx.Image)

	// step 3. run the image (it *should* just echo in response)
	if err := shell.RunAndLog("docker", "run", "-d", "-p", "9800:9800", "--name", "plumber-test", ctx.Image); err != nil {
		t.Errorf("Got an error during docker run: '%v'", err)
	}
	defer shell.RunAndLog("docker", "rm", "-f", "plumber-test")
	// wait a bit for the container to come up
	time.Sleep(1 * time.Second)

	// step 4. send some JSON and check for echos
	// first, find the IP to connect to

	// get the DOCKER_HOST environment variable; if not defined, use
	// docker to find it
	hostIp := os.Getenv("DOCKER_HOST")
	if hostIp == "" {
		cmd := exec.Command("docker", "inspect", "--format='{{.NetworkSettings.Gateway}}'", "plumber-test")
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

	// second, send over some JSON and verify result
	resp, err := http.Post(fmt.Sprintf("http://%s:9800", hostIp), "application/json", bytes.NewBufferString(`{"foo": 3}`))
	if err != nil {
		t.Error(err)
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	result := buf.String()
	if result != `{"foo": 3}` {
		t.Errorf("Got '%s'; did not get expected response", result)
	}
}
