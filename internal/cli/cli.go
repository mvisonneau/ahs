package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/mvisonneau/ahs/internal/cmd"
	"github.com/urfave/cli/v2"
)

// Run handles the instanciation of the CLI application.
func Run(version string, args []string) {
	err := NewApp(version, time.Now()).Run(args)
	if err != nil {
		fmt.Println(err) //nolint
		os.Exit(1)
	}
}

// NewApp configures the CLI application.
func NewApp(version string, start time.Time) (app *cli.App) {
	app = cli.NewApp()
	app.Name = "ahs"
	app.Version = version
	app.Usage = "Set the hostname of an EC2 instance based on a tag value and the instance-id"
	app.EnableBashCompletion = true

	app.Flags = cli.FlagsByName{
		&cli.BoolFlag{
			Name:    "dry-run",
			EnvVars: []string{"AHS_DRY_RUN"},
			Usage:   "only display what would have been done",
		},
		&cli.StringFlag{
			Name:    "input-tag",
			EnvVars: []string{"AHS_INPUT_TAG"},
			Usage:   "`tag` to use as input to determine the hostname",
			Value:   "Name",
		},
		&cli.StringFlag{
			Name:    "log-level",
			EnvVars: []string{"AHS_LOG_LEVEL"},
			Usage:   "log `level` (debug,info,warn,fatal,panic)",
			Value:   "info",
		},
		&cli.StringFlag{
			Name:    "log-format",
			EnvVars: []string{"AHS_LOG_FORMAT"},
			Usage:   "log `format` (json,text)",
			Value:   "text",
		},
		&cli.StringFlag{
			Name:    "output-tag",
			EnvVars: []string{"AHS_OUTPUT_TAG"},
			Usage:   "`tag` to update with the computed hostname",
			Value:   "Name",
		},
		&cli.BoolFlag{
			Name:    "persist-hostname",
			EnvVars: []string{"AHS_PERSIST_HOSTNAME"},
			Usage:   "set /etc/hostname with generated hostname",
		},
		&cli.BoolFlag{
			Name:    "persist-hosts",
			EnvVars: []string{"AHS_PERSIST_HOSTS"},
			Usage:   "assign generated hostname to 127.0.0.1 in /etc/hosts",
		},
		&cli.StringFlag{
			Name:    "separator",
			EnvVars: []string{"AHS_SEPARATOR"},
			Usage:   "`separator` to use between tag and id",
			Value:   "-",
		},
	}

	app.Commands = cli.CommandsByName{
		{
			Name:      "instance-id",
			Usage:     "compute a hostname by appending the instance-id to a prefixed/base string",
			ArgsUsage: " ",
			Flags: cli.FlagsByName{
				&cli.IntFlag{
					Name:    "length",
					EnvVars: []string{"AHS_INSTANCE_ID_LENGTH"},
					Usage:   "length of the id to keep in the hostname",
					Value:   5,
				},
			},
			Action: cmd.ExecWrapper(cmd.Run),
		},
		{
			Name:      "sequential",
			Usage:     "compute a sequential hostname based on the number of instances belonging to the same group",
			ArgsUsage: " ",
			Flags: cli.FlagsByName{
				&cli.StringFlag{
					Name:    "instance-sequential-id-tag",
					EnvVars: []string{"AHS_INSTANCE_SEQUENTIAL_ID_TAG"},
					Usage:   "tag to which output the computed instance-sequential-id",
					Value:   "ahs:instance-id",
				},
				&cli.StringFlag{
					Name:    "instance-group-tag",
					EnvVars: []string{"AHS_INSTANCE_GROUP_TAG"},
					Usage:   "tag to use in order to determine which group the instance belongs to",
					Value:   "ahs:instance-group",
				},
				&cli.BoolFlag{
					Name:    "respect-azs",
					EnvVars: []string{"AHS_RESPECT_AZS"},
					Usage:   "if instances are provisioned through an ASG, setting this flag it will get the sequential-ids associated to respective azs", //nolint
				},
				&cli.StringFlag{
					Name:    "valid-instance-states",
					EnvVars: []string{"AHS_VALID_INSTANCE_STATES"},
					Usage:   "comma-delimited list selecting which instance states (running, stopped, etc.) are valid when filtering instances to take into account for assigning sequential ids",
					Value:   "running",
				},
			},
			Action: cmd.ExecWrapper(cmd.Run),
		},
	}

	app.Metadata = map[string]interface{}{
		"startTime": start,
	}

	return
}
