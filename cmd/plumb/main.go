package main

import (
	"errors"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/qadium/plumb"
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
	app.Name = "plumb"
	app.Usage = "a command line tool for managing information discovery"
	app.Version = versionString
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "server, s",
			Value:  "/var/run/plumb.sock",
			Usage:  "location of plumb server socket",
			EnvVar: "LINK_SERVER",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:   "create",
			Usage:  "create a pipeline managed by plumb",
			Before: createRequiredArgCheck(exactly(1), "Please provide a pipeline name."),
			Action: func(c *cli.Context) {
				path := c.Args().First()
				if err := plumb.Create(path); err != nil {
					panic(err)
				}
			},
		},
		{
			Name:   "add",
			Usage:  "add a plumb-enabled bundle to a pipeline",
			Before: createRequiredArgCheck(atLeast(2), "Please provide both a pipeline name and a bundle path."),
			Action: func(c *cli.Context) {
				pipeline := c.Args()[0]
				bundles := c.Args()[1:]
				if err := plumb.Add(pipeline, bundles...); err != nil {
					panic(err)
				}
			},
		},
		{
			Name:   "start",
			Usage:  "start a pipeline managed by plumb",
			Before: createRequiredArgCheck(exactly(1), "Please provide a pipeline name."),
			Action: func(c *cli.Context) {
				pipeline := c.Args().First()
				if err := plumb.Start(pipeline); err != nil {
					panic(err)
				}
			},
		},
		{
			Name:   "bundle",
			Usage:  "bundle a node for use in a pipeline managed by plumb",
			Before: createRequiredArgCheck(exactly(1), "Please provide a bundle path."),
			Action: func(c *cli.Context) {
				path := c.Args().First()
				if err := plumb.Bundle(path); err != nil {
					panic(err)
				}
			},
		},
		{
			Name:  "version",
			Usage: "more detailed version information for plumb",
			Action: func(c *cli.Context) {
				fmt.Println("plumb version:", versionString)
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
