package cli

import (
	"fmt"
	"github.com/qadium/plumber/bindata"
	"github.com/qadium/plumber/shell"
	"io/ioutil"
	"log"
	"os/user"
	"os"
	// "gopkg.in/yaml.v2"
)

func writeAsset(asset string, directory string) error {
	log.Printf(" |     Writing '%s'", asset)
	data, err := bindata.Asset(fmt.Sprintf("manager/%s", asset))
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(fmt.Sprintf("%s/%s", directory, asset), data, 0644); err != nil {
		return err
	}
	log.Printf("       Done.")
	return nil
}

func Bootstrap() error {
	// use docker to compile the manager and copy the binary into
	// another docker container
	//
	// tag this new container
	log.Printf("==> Bootstraping plumb.")
	defer log.Printf("<== Bootstrap complete.")

	log.Printf(" |  Creating temp directory.")
	usr, err := user.Current()
	if err != nil {
		return err
	}
	directory := fmt.Sprintf("%s/.plumber-bootstrap", usr.HomeDir)

	if err := os.MkdirAll(directory, 0755); err != nil {
		return err
	}
	defer func() {
		if err := os.RemoveAll(directory); err != nil {
			panic(err)
		}
	}()
	log.Printf("    Temp directory created at '%s'", directory)

	log.Printf(" |  Writing manager source files.")
	if err := writeAsset("manager.go", directory); err != nil {
		return err
	}
	if err := writeAsset("Dockerfile", directory); err != nil {
		return err
	}
	log.Printf("    Done")

	if err := shell.RunAndLog("docker",
		"run",
		"--rm",
		"-v",
		"/var/run/docker.sock:/var/run/docker.sock",
		"-v",
		fmt.Sprintf("%s:/src", directory),
		"centurylink/golang-builder",
		"plumber/manager"); err != nil {
		return err
	}

	return nil
}
