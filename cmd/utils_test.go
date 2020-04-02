package cmd

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExit(t *testing.T) {
	err := exit(20, fmt.Errorf("test"))
	assert.Equal(t, err.Error(), "")
	assert.Equal(t, err.ExitCode(), 20)
}

func TestAnalyzeEC2APIError(t *testing.T) {
	assert.Equal(t, "", analyzeEC2APIError(nil))

	err := errors.New("test")
	assert.Equal(t, err.Error(), analyzeEC2APIError(err))
}
