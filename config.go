package main

var cfg Config

// Config : A type that handles the configuration of the app
type Config struct {
	DryRun   bool
	IDLength int
	InputTag string

	Log struct {
		Level  string
		Format string
	}

	OutputTag string
	Separator string
}
