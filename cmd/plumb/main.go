package main

import (
	"errors"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/qadium/plumb"
	"os"
)

func createRequiredArgCheck(num int, message string) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		if len(c.Args()) != num {
			fmt.Println(message)
			return errors.New(message)
		}
		return nil
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
			Before: createRequiredArgCheck(1, "Please provide a pipeline name."),
			Action: func(c *cli.Context) {
				name := c.Args().First()
				fmt.Println("CREATE", name)
			},
		},
		{
			Name:   "submit",
			Usage:  "submit a plumb created bundle to a pipeline",
			Before: createRequiredArgCheck(2, "Please provide both a pipeline name and a bundle path."),
			Action: func(c *cli.Context) {
				name := c.Args()[0]
				pipeline := c.Args()[1]
				fmt.Println("SUBMIT", name, pipeline)
			},
		},
		{
			Name:   "start",
			Usage:  "start a pipeline managed by plumb",
			Before: createRequiredArgCheck(1, "Please provide a pipeline name."),
			Action: func(c *cli.Context) {
				name := c.Args().First()
				fmt.Println("START", name)
			},
		},
		{
			Name:   "bundle",
			Usage:  "bundle a node for use in a pipeline managed by plumb",
			Before: createRequiredArgCheck(1, "Please provide a bundle path."),
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
