package cli

import (
	"fmt"
	"log"
	"github.com/qadium/plumb/shell"
	// "gopkg.in/yaml.v2"
	// "io/ioutil"
)

func Bootstrap(commit string) error {
	log.Printf("==> Bootstraping plumb.")
	defer log.Printf("<== Bootstrap complete.")

	if err := shell.RunAndLog("docker",
		"run",
		"--rm",
		"--entrypoint=\"/bin/bash\"",
		"-v",
		"/var/run/docker.sock:/var/run/docker.sock",
		"centurylink/golang-builder",
		"-c",
		fmt.Sprintf(`git clone https://github.com/qadium/plumb /plumb \
&& cd /plumb \
&& git checkout -b build %s \
&& cp -r /plumb/manager/* /src \
&& cd /src \
&& /build.sh`, commit)); err != nil {
		return err
	}

	// use docker to compile the manager and copy the binary into
	// another docker container
	//
	// tag this new container

	return nil
}
