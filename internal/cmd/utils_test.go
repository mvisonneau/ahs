package cmd

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

func TestAnalyzeEC2APIError(t *testing.T) {
	assert.Equal(t, "", analyzeEC2APIError(nil))

	err := errors.New("test")
	assert.Equal(t, err.Error(), analyzeEC2APIError(err))
}

func TestExit(t *testing.T) {
	err := exit(20, fmt.Errorf("test"))
	assert.Equal(t, "", err.Error())
	assert.Equal(t, 20, err.ExitCode())
}

func TestExecWrapper(t *testing.T) {
	function := func(ctx *cli.Context) (int, error) {
		return 0, nil
	}
	assert.Equal(t, exit(function(&cli.Context{})), ExecWrapper(function)(&cli.Context{}))
}
