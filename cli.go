package main

import (
	"github.com/urfave/cli"
)

var version = "<devel>"

// runCli : Generates cli configuration for the application
func runCli() (c *cli.App) {
	c = cli.NewApp()
	c.Name = "ahs"
	c.Version = version
	c.Usage = "Set the hostname of an EC2 instance based on a tag value and the instance-id"
	c.EnableBashCompletion = true

	c.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "dry-run",
			EnvVar:      "AHS_DRY_RUN",
			Usage:       "only display what would have been done",
			Destination: &cfg.DryRun,
		},
		cli.IntFlag{
			Name:        "id-length",
			EnvVar:      "AHS_ID_LENGTH",
			Usage:       "length of the id to keep in the hostname",
			Value:       5,
			Destination: &cfg.IDLength,
		},
		cli.StringFlag{
			Name:        "input-tag",
			EnvVar:      "AHS_TAG_NAME_INPUT",
			Usage:       "tag to use as input to determine the hostname",
			Value:       "Name",
			Destination: &cfg.InputTag,
		},
		cli.StringFlag{
			Name:        "log-level",
			EnvVar:      "AHS_LOG_LEVEL",
			Usage:       "log level (debug,info,warn,fatal,panic)",
			Value:       "info",
			Destination: &cfg.Log.Level,
		},
		cli.StringFlag{
			Name:        "log-format",
			EnvVar:      "AHS_LOG_FORMAT",
			Usage:       "log format (json,text)",
			Value:       "text",
			Destination: &cfg.Log.Format,
		},
		cli.StringFlag{
			Name:        "output-tag",
			EnvVar:      "AHS_TAG_NAME_OUTPUT",
			Usage:       "tag to update with the computed hostname",
			Value:       "Name",
			Destination: &cfg.OutputTag,
		},
		cli.StringFlag{
			Name:        "separator",
			EnvVar:      "AHS_SEPARATOR",
			Usage:       "separator to use between tag and id",
			Value:       "-",
			Destination: &cfg.Separator,
		},
	}

	c.Commands = []cli.Command{
		{
			Name:      "run",
			Usage:     "replace the hostname with found/computed values",
			ArgsUsage: " ",
			Action:    run,
		},
	}

	return
}
