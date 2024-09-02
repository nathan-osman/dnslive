package main

import (
	"github.com/nathan-osman/dnslive/server"
	"github.com/nathan-osman/gosvc"
	"github.com/urfave/cli/v2"
)

func serverCommand() *cli.Command {
	var (
		a = &gosvc.Application{
			Name:            "dnslive-server",
			Description:     "DNS server for dynamic IP addresses",
			Args:            []string{"server"},
			RequiresNetwork: true,
		}
		p = a.Platform()
	)
	return &cli.Command{
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
		Subcommands: gosvc.Commands(p),
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
			return p.Run()
		},
	}
}
