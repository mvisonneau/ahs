package cmd

import (
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/mvisonneau/go-helpers/logger"
	"github.com/urfave/cli"

	log "github.com/sirupsen/logrus"
)

var start time.Time

func configure(ctx *cli.Context) error {
	start = ctx.App.Metadata["startTime"].(time.Time)

	lc := &logger.Config{
		Level:  ctx.GlobalString("log-level"),
		Format: ctx.GlobalString("log-format"),
	}

	return lc.Configure()
}

func exit(exitCode int, err error) *cli.ExitError {
	defer log.Debugf("Executed in %s, exiting..", time.Since(start))
	if err != nil {
		log.Error(analyzeEC2APIError(err))
	}

	return cli.NewExitError("", exitCode)
}

func analyzeEC2APIError(err error) string {
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			return aerr.Error()
		}
		return err.Error()
	}
	return ""
}

// ExecWrapper gracefully logs and exits our `run` functions
func ExecWrapper(f func(ctx *cli.Context) (int, error)) func(*cli.Context) error {
	return func(ctx *cli.Context) error {
		return exit(f(ctx))
	}
}
