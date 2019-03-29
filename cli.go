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
			Name:   "dry-run",
			EnvVar: "AHS_DRY_RUN",
			Usage:  "only display what would have been done",
		},
		cli.StringFlag{
			Name:   "input-tag",
			EnvVar: "AHS_INPUT_TAG",
			Usage:  "`tag` to use as input to determine the hostname",
			Value:  "Name",
		},
		cli.StringFlag{
			Name:   "log-level",
			EnvVar: "AHS_LOG_LEVEL",
			Usage:  "log `level` (debug,info,warn,fatal,panic)",
			Value:  "info",
		},
		cli.StringFlag{
			Name:   "log-format",
			EnvVar: "AHS_LOG_FORMAT",
			Usage:  "log `format` (json,text)",
			Value:  "text",
		},
		cli.StringFlag{
			Name:   "output-tag",
			EnvVar: "AHS_OUTPUT_TAG",
			Usage:  "`tag` to update with the computed hostname",
			Value:  "Name",
		},
		cli.BoolFlag{
			Name:   "persist-hostname",
			EnvVar: "AHS_PERSIST_HOSTNAME",
			Usage:  "set /etc/hostname with generated hostname",
		},
		cli.BoolFlag{
			Name:   "persist-hosts",
			EnvVar: "AHS_PERSIST_HOSTS",
			Usage:  "assign generated hostname to 127.0.0.1 in /etc/hosts",
		},
		cli.StringFlag{
			Name:   "separator",
			EnvVar: "AHS_SEPARATOR",
			Usage:  "`separator` to use between tag and id",
			Value:  "-",
		},
	}

	c.Commands = []cli.Command{
		{
			Name:      "instance-id",
			Usage:     "compute a hostname by appending the instance-id to a prefixed/base string",
			ArgsUsage: " ",
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:   "length",
					EnvVar: "AHS_INSTANCE_ID_LENGTH",
					Usage:  "length of the id to keep in the hostname",
					Value:  5,
				},
			},
			Action: run,
		},
		{
			Name:      "sequential",
			Usage:     "compute a sequential hostname based on the number of instances belonging to the same group",
			ArgsUsage: " ",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "instance-sequential-id-tag",
					EnvVar: "AHS_INSTANCE_SEQUENTIAL_ID_TAG",
					Usage:  "tag to which output the computed instance-sequential-id",
					Value:  "ahs:instance-id",
				},
				cli.StringFlag{
					Name:   "instance-group-tag",
					EnvVar: "AHS_INSTANCE_GROUP_TAG",
					Usage:  "tag to use in order to determine which group the instance belongs to",
					Value:  "ahs:instance-group",
				},
				cli.BoolFlag{
					Name:   "respect-azs",
					EnvVar: "AHS_RESPECT_AZS",
					Usage:  "if instances are provisioned through an ASG, setting this flag it will get the sequential-ids associated to respective azs",
				},
			},
			Action: run,
		},
	}

	return
}
