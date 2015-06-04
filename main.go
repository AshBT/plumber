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
	versionString := versionString()
	app := cli.NewApp()
	app.Name = "plumber"
	app.Usage = "a command line tool for managing distributed data pipelines"
	app.Version = versionString
	// app.Flags = []cli.Flag{
	// 	cli.StringFlag{
	// 		Name:   "server, s",
	// 		Value:  "/var/run/plumber.sock",
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
				if err := plumber.Add(pipeline, bundles...); err != nil {
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
				if err := plumber.Create(path); err != nil {
					panic(err)
				}
			},
		},
		{
			Name:  "bootstrap",
			Usage: "bootstrap local setup for use with plumber",
			Description: `The bootstrap command builds the latest manager for use with plumber.
This packages the manager into a minimal container for use on localhost.

When running the pipeline on Google Cloud, the manager container is
pushed to your project's private repository.`,
			Action: func(c *cli.Context) {
				if err := plumber.Bootstrap(); err != nil {
					panic(err)
				}
			},
		},
		{
			Name:   "start",
			Usage:  "start a pipeline managed by plumber",
			Flags:  []cli.Flag {
				cli.StringFlag{
					Name: "gce",
					Value: "",
					Usage: "Google Cloud project ID",
				},
			},
			Before: createRequiredArgCheck(exactly(1), "Please provide a pipeline name."),
			Action: func(c *cli.Context) {
				pipeline := c.Args().First()
				gce := c.String("gce")
				if err := plumber.Start(pipeline, gce, versionString, GitCommit); err != nil {
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
				if err := plumber.Bundle(path); err != nil {
					panic(err)
				}
			},
		},
		{
			Name:  "version",
			Usage: "more detailed version information for plumber",
			Action: func(c *cli.Context) {
				fmt.Println("plumber version:", versionString)
				fmt.Println("git commit:", GitCommit)
			},
		},
	}
	app.Run(os.Args)
}

func versionString() string {
	versionString := Version
	if VersionPrerelease != "" {
		versionString += "-" + VersionPrerelease
	}
	return versionString
}
