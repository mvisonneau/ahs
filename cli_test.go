package main

import (
	"testing"
)

func TestRunCli(t *testing.T) {
	c := runCli()
	if c.Name != "ahs" {
		t.Fatalf("Expected c.Name to be ahs, got '%v'", c.Name)
	}
}
