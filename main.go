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
// Plumber is a cli tool for creating distributed data processing
// pipelines. The main package only contains the driver for the cli
// tool.
package main

import (
	"errors"
	"fmt"
	"github.com/codegangsta/cli"
	plumber "github.com/qadium/plumber/cli"
	"os"
)

func createRequiredArgCheck(check func(args cli.Args) bool, message string) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		if !check(c.Args()) {
			fmt.Println(message)
			return errors.New(message)
		}
		return nil
	}
}

func exactly(num int) func(args cli.Args) bool {
	return func(args cli.Args) bool {
		return len(args) == num
	}
}

func atLeast(num int) func(args cli.Args) bool {
	return func(args cli.Args) bool {
		return len(args) >= num
	}
}

func main() {
	plumberCtx, err := plumber.NewDefaultContext()
	if err != nil {
		panic(err)
	}
	app := cli.NewApp()
	app.Name = "plumber"
	app.Usage = "a command line tool for managing distributed data pipelines"
	app.Version = plumberCtx.Version
	// app.Flags = []cli.Flag{
	// 	cli.StringFlag{
	// 		Name:   "server, s",
	// 		Value:  "/var/run/plumberCtx.sock",
	// 		Usage:  "location of plumber server socket",
	// 		EnvVar: "LINK_SERVER",
	// 	},
	// }
	app.Commands = []cli.Command{
		{
			Name:   "add",
			Usage:  "add a plumber-enabled bundle to a pipeline",
			Before: createRequiredArgCheck(atLeast(2), "Please provide both a pipeline name and a bundle path."),
			Action: func(c *cli.Context) {
				pipeline := c.Args()[0]
				bundles := c.Args()[1:]
				if err := plumberCtx.Add(pipeline, bundles...); err != nil {
					panic(err)
				}
			},
		},
		{
			Name:   "create",
			Usage:  "create a pipeline managed by plumber",
			Before: createRequiredArgCheck(exactly(1), "Please provide a pipeline name."),
			Action: func(c *cli.Context) {
				path := c.Args().First()
				if err := plumberCtx.Create(path); err != nil {
					panic(err)
				}
			},
		},
		{
			Name:  "bootstrap",
			Usage: "bootstrap local setup for use with plumber",
			Description: `The bootstrap command builds the latest manager for use with plumberCtx.
This packages the manager into a minimal container for use on localhost.

When running the pipeline on Google Cloud, the manager container is
pushed to your project's private repository.`,
			Action: func(c *cli.Context) {
				if err := plumberCtx.Bootstrap(); err != nil {
					panic(err)
				}
			},
		},
		{
			Name:  "start",
			Usage: "start a pipeline managed by plumber",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "gce",
					Value: "",
					Usage: "Google Cloud project ID",
				},
			},
			Before: createRequiredArgCheck(exactly(1), "Please provide a pipeline name."),
			Action: func(c *cli.Context) {
				pipeline := c.Args().First()
				gce := c.String("gce")
				if err := plumberCtx.Start(pipeline, gce); err != nil {
					panic(err)
				}
			},
		},
		{
			Name:   "bundle",
			Usage:  "bundle a node for use in a pipeline managed by plumber",
			Before: createRequiredArgCheck(exactly(1), "Please provide a bundle path."),
			Action: func(c *cli.Context) {
				path := c.Args().First()
				if err := plumberCtx.Bundle(path); err != nil {
					panic(err)
				}
			},
		},
		{
			Name:  "version",
			Usage: "more detailed version information for plumber",
			Action: func(c *cli.Context) {
				fmt.Println("plumber version:", plumberCtx.Version)
				fmt.Println("git commit:", plumberCtx.GitCommit)
			},
		},
	}
	app.Run(os.Args)
}
