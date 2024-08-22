package main

import (
	"github.com/urfave/cli/v2"
)

var clientCommand = &cli.Command{
	Name:  "client",
	Usage: "run the application in client mode",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "config",
			EnvVars: []string{"CONFIG"},
			Usage:   "path to a configuration file",
		},
	},
}
