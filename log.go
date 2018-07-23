package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

// Logger config
type Logger struct {
	Level  string
	Format string
}

// Configure the logger
func (l *Logger) Configure() error {
	parsedLevel, err := log.ParseLevel(l.Level)
	if err != nil {
		return err
	}
	log.SetLevel(parsedLevel)

	formatter := &log.TextFormatter{
		FullTimestamp: true,
	}
	log.SetFormatter(formatter)

	switch l.Format {
	case "text":
	case "json":
		log.SetFormatter(&log.JSONFormatter{})
	default:
		return fmt.Errorf("Invalid log format '%s'", l.Format)
	}

	log.SetOutput(os.Stdout)

	return nil
}
