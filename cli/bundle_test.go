package cli_test

import (
	"fmt"
	"github.com/qadium/plumber/shell"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
	//"log"
	"bytes"
)

// put bundle under test by inputting yaml from string and checking
// that output Dockerfile is expected
//
// would be nice to run the container and curl it, but that is also
// checked manually
//
// we manually check that output Dockerfile can be built
//
// so we're just checking that *strings* match, but not that the
// functionality is as expected
//
// we probably shouldn't test the wrappers anyway, since they should be
// tested in a separate repo (so we can have language agnosticism)

// We'll check that the Bundle command actually works

// TODO: this test code should live in the language-dependent repository
// and our test should *load* those tests
//
// this test echoes the input
const bundleTestPython = `
def run(in_dict):
	in_dict['b'] = "echo {}".format(in_dict['a'])
	return in_dict
`

func TestBundle(t *testing.T) {
	//log.SetOutput(ioutil.Discard)

	ctx, tempDir := NewTestContext(t)
	defer cleanTestDir(t, tempDir)

	// use the `goodBundle` from config_test.go
	filename := fmt.Sprintf("%s/.plumb.yml", tempDir)
	if err := ioutil.WriteFile(filename, []byte(optBundle), 0644); err != nil {
		t.Errorf("Could not write .plumb.yml; got error '%v'", err)
	}

	// write out the python test code
	filename = fmt.Sprintf("%s/foobar.py", tempDir)
	if err := ioutil.WriteFile(filename, []byte(bundleTestPython), 0644); err != nil {
		t.Errorf("Could not write foobar.py; got error '%v'", err)
	}

	// write out the empty requirements.txt file
	filename = fmt.Sprintf("%s/requirements.txt", tempDir)
	if err := ioutil.WriteFile(filename, []byte{}, 0644); err != nil {
		t.Errorf("Could not write requirements.txt; got error '%v'", err)
	}

	// bundle the tempDir
	err := ctx.Bundle(tempDir)
	if err != nil {
		t.Errorf("Got an unxpected error while bundling, '%v'", err)
	}
	defer shell.RunAndLog("docker", "rmi", ctx.GetImage("foobar"))

	// run the container and check that it increments the input
	if err := shell.RunAndLog("docker", "run", "-d", "-p", "9800:9800", "--name", "foobar", ctx.GetImage("foobar")); err != nil {
		t.Errorf("Got an error during docker run: '%v'", err)
	}
	defer shell.RunAndLog("docker", "rm", "-f", "foobar")
	// wait a bit for the container to come up
	time.Sleep(1 * time.Second)

	// step 4. send some JSON and check for echos
	// first, find the IP to connect to
	hostIp := getImageIp(t, "foobar")

	// second, send over some JSON and verify result
	resp, err := http.Post(fmt.Sprintf("http://%s:9800", hostIp), "application/json", bytes.NewBufferString(`{"a": "trusty"}`))
	if err != nil {
		t.Error(err)
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	result := buf.String()
	if result != `{"a": "trusty", "b": "echo trusty"}` {
		t.Errorf("Got '%s'; did not get expected response", result)
	}
}
