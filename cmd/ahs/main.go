package main

import (
	"os"

	"github.com/mvisonneau/ahs/internal/cli"
)

var version = ""

func main() {
	cli.Run(version, os.Args)
}
