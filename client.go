package main

import (
	"time"

	"github.com/nathan-osman/dnslive/client"
	"github.com/urfave/cli/v2"
)

var clientCommand = &cli.Command{
	Name:  "client",
	Usage: "run the application in client mode",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "config",
			Usage:    "path to a configuration file",
			Required: true,
			EnvVars:  []string{"CONFIG"},
		},
	},
	Action: func(ctx *cli.Context) error {
		v := client.Config{
			Interval: 1 * time.Hour,
		}
		if err := readConfigFile(ctx.String("config"), &v); err != nil {
			return err
		}
		c, err := client.New(&v)
		if err != nil {
			return err
		}
		defer c.Close()
		return run()
	},
}
