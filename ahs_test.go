package main

import (
	"errors"
	"io/ioutil"
	"testing"

	"github.com/txn2/txeh"
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

func TestComputeHostnameWithInstanceID(t *testing.T) {
	hostname, err := computeHostnameWithInstanceID("myhostname", "-", "i-123456789", 5)
	if err != nil {
		t.Fatalf("Shouldn't have returned any error, got : '%s'", err.Error())
	}

	if hostname != "myhostname-12345" {
		t.Fatalf("Should have retreived myhostname-12345, got '%s'", hostname)
	}

	hostname, err = computeHostnameWithInstanceID("myhostname-12345", "-", "i-123456789", 5)
	if err != nil {
		t.Fatalf("Shouldn't have returned any error, got : '%s'", err.Error())
	}

	if hostname != "myhostname-12345" {
		t.Fatalf("Should have retreived myhostname-12345, got '%s'", hostname)
	}
}

func TestValidComputeRegionFromAZ(t *testing.T) {
	region, err := computeRegionFromAZ("eu-west-1a")
	if err != nil {
		t.Fatalf("Shouldn't have returned any error, got : '%s'", err.Error())
	}

	if region != "eu-west-1" {
		t.Fatalf("Should have retreived eu-west-1, got '%s'", region)
	}
}

func TestInvalidComputeRegionFromAZ(t *testing.T) {
	_, err := computeRegionFromAZ("foo")
	if err == nil {
		t.Fatal("Should have thrown an error, got nil")
	}
}

func TestUpdateHostnameFile(t *testing.T) {
	if err := updateHostnameFile("myhostname"); err != nil {
		t.Fatalf("Shouldn't have returned any error, got : '%s'", err.Error())
	}

	content, err := ioutil.ReadFile("/etc/hostname")
	if err != nil {
		t.Fatalf("Should have been able to read /etc/hostname without error, got : '%s'", err.Error())
	}

	if string(content) != "myhostname\n" {
		t.Fatalf("/etc/hostname content should be equal to 'myhostname' but is '%s'", content)
	}
}

func TestUpdateHostsFile(t *testing.T) {
	if err := updateHostsFile("myhostname"); err != nil {
		t.Fatalf("Shouldn't have returned any error, got : '%s'", err.Error())
	}

	hosts, err := txeh.NewHostsDefault()
	if err != nil {
		t.Fatalf("Shouldn't have returned any error, got : '%s'", err.Error())
	}

	found, address, _ := hosts.HostAddressLookup("myhostname")
	if !found {
		t.Fatalf("'myhostname' host could not be found in /etc/hosts")
	}

	if address != "127.0.0.1" {
		t.Fatalf("'myhostname' address is not equal to 127.0.0.1, got '%s'", address)
	}
}

func TestExit(t *testing.T) {
	err := errors.New("test")

	if exit(err) != err {
		t.Fatalf("Error should be equal to the orignal one")
	}
}
