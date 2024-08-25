package main

import (
	"github.com/nathan-osman/dnslive/server"
	"github.com/urfave/cli/v2"
)

var serverCommand = &cli.Command{
	Name:  "server",
	Usage: "run the application in server mode",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "config",
			Usage:    "path to a configuration file",
			Required: true,
			EnvVars:  []string{"CONFIG"},
		},
	},
	Action: func(ctx *cli.Context) error {
		v := server.Config{
			HttpServerAddr: "0.0.0.0:443",
			DnsServerAddr:  "0.0.0.0:53",
			PersistentFile: "entries.json",
		}
		if err := readConfigFile(ctx.String("config"), &v); err != nil {
			return err
		}
		s, err := server.New(&v)
		if err != nil {
			return err
		}
		defer s.Close()
		waitForInterrupt()
		return nil
	},
}
