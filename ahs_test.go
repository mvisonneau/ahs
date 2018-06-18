package main

import (
	"errors"
  "testing"
)

func TestAnalyzeEC2APIErrors(t *testing.T) {
  if analyzeEC2APIErrors(nil) != "" {
    t.Fatalf("Expected an empty string, got '%s'", analyzeEC2APIErrors(nil))
  }

  err := errors.New("test")
  if analyzeEC2APIErrors(err) != err.Error() {
    t.Fatalf("Expected to return error content, got '%s'", analyzeEC2APIErrors(err))
  }
}

func TestComputeHostname(t *testing.T) {
  hostname := computeHostname("myhostname", "-", "i-123456789", 5 )
  if hostname != "myhostname-12345" {
    t.Fatalf("Should have retreived myhostname-12345, got '%s'", hostname)
  }
}

func TestComputeRegionFromAZ(t *testing.T) {
  region := computeRegionFromAZ("eu-west-1a")
  if region != "eu-west-1" {
    t.Fatalf("Should have retreived eu-west-1, got '%s'", region)
  }
}

func TestExit(t *testing.T) {
  err := errors.New("test")

  if exit(err) != err {
    t.Fatalf("Error should be equal to the orignal one")
  }
}
