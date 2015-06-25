/**
 * Copyright 2015 Qadium, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package cli

import (
	"fmt"
	"github.com/qadium/plumber/bindata"
	"github.com/qadium/plumber/shell"
	"io/ioutil"
	"log"
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

func (ctx *Context) Bootstrap() error {
	// use docker to compile the manager and copy the binary into
	// another docker container
	//
	// tag this new container
	log.Printf("==> Bootstraping plumb.")
	defer log.Printf("<== Bootstrap complete.")

	log.Printf(" |  Creating temp directory.")
	if err := os.MkdirAll(ctx.BootstrapDir, 0755); err != nil {
		return err
	}
	defer func() {
		if err := os.RemoveAll(ctx.BootstrapDir); err != nil {
			panic(err)
		}
	}()
	log.Printf("    Temp directory created at '%s'", ctx.BootstrapDir)

	log.Printf(" |  Writing manager source files.")
	if err := writeAsset("manager.go", ctx.BootstrapDir); err != nil {
		return err
	}
	if err := writeAsset("Dockerfile", ctx.BootstrapDir); err != nil {
		return err
	}
	if err := writeAsset("README.md", ctx.BootstrapDir); err != nil {
		return err
	}
	if err := writeAsset("manager_test.go", ctx.BootstrapDir); err != nil {
		return err
	}
	log.Printf("    Done")

	if err := shell.RunAndLog(ctx.DockerCmd,
		"run",
		"--rm",
		"-v",
		"/var/run/docker.sock:/var/run/docker.sock",
		"-v",
		fmt.Sprintf("%s:/src", ctx.BootstrapDir),
		"centurylink/golang-builder",
		ctx.GetManagerImage()); err != nil {
		return err
	}

	return nil
}
