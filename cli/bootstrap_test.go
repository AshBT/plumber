package cli_test

import (
	"bytes"
	"fmt"
	"github.com/qadium/plumber/shell"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"
)

// this is a *functional test*
// we check that boostrap works by actually running the boostrap command
// and checking that the container is built and runs
func TestBootstrap(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	ctx, tempDir := NewTestContext(t)
	defer cleanTestDir(t, tempDir)

	// step 1. remove any image named plumber/test-manager from the current
	// set of docker images (ignore any errors)
	_ = shell.RunAndLog("docker", "rmi", ctx.GetManagerImage())

	// step 2. invoke Bootstrap for building ctx.GetManagerImage()
	if err := ctx.Bootstrap(); err != nil {
		t.Errorf("Got an error during bootstrap: '%v'", err)
	}
	defer shell.RunAndLog("docker", "rmi", ctx.GetManagerImage())

	// step 3. run the image (it *should* just echo in response)
	if err := shell.RunAndLog("docker", "run", "-d", "-p", "9800:9800", "--name", "plumber-test", ctx.GetManagerImage()); err != nil {
		t.Errorf("Got an error during docker run: '%v'", err)
	}
	defer shell.RunAndLog("docker", "rm", "-f", "plumber-test")
	// wait a bit for the container to come up
	time.Sleep(1 * time.Second)

	// step 4. send some JSON and check for echos
	hostIp := getImageIp(t, "plumber-test")

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
