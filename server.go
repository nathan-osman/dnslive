package main

import (
	"github.com/urfave/cli/v2"
)

var serverCommand = &cli.Command{
	Name:  "server",
	Usage: "run the application in server mode",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "config",
			EnvVars: []string{"CONFIG"},
			Usage:   "path to a configuration file",
		},
	},
}
