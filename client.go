package main

import (
	"time"

	"github.com/nathan-osman/dnslive/client"
	"github.com/nathan-osman/gosvc"
	"github.com/urfave/cli/v2"
)

func clientCommand() *cli.Command {
	var (
		a = &gosvc.Application{
			Name:            "dnslive-client",
			Description:     "DNS client for dynamic IP addresses",
			Args:            []string{"client"},
			RequiresNetwork: true,
		}
		p = a.Platform()
	)
	return &cli.Command{
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
		Subcommands: gosvc.Commands(p),
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
			return p.Run()
		},
	}
}
