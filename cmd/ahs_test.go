package cmd

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/txn2/txeh"
)

func TestComputeHostnameWithInstanceID(t *testing.T) {
	hostname, err := computeHostnameWithInstanceID("myhostname", "-", "i-123456789", 5)
	assert.Nil(t, err)
	assert.Equal(t, "myhostname-12345", hostname)

	hostname, err = computeHostnameWithInstanceID("myhostname-12345", "-", "i-123456789", 5)
	assert.Nil(t, err)
	assert.Equal(t, "myhostname-12345", hostname)
}

func TestValidComputeRegionFromAZ(t *testing.T) {
	region, err := computeRegionFromAZ("eu-west-1a")
	assert.Nil(t, err)
	assert.Equal(t, "eu-west-1", region)
}

func TestInvalidComputeRegionFromAZ(t *testing.T) {
	_, err := computeRegionFromAZ("foo")
	assert.NotNil(t, err)
}

func TestUpdateHostnameFile(t *testing.T) {
	err := updateHostnameFile("myhostname")
	assert.Nil(t, err)

	content, err := ioutil.ReadFile("/etc/hostname")
	assert.Nil(t, err)
	assert.Equal(t, "myhostname\n", string(content))
}

func TestUpdateHostsFile(t *testing.T) {
	err := updateHostsFile("myhostname")
	assert.Nil(t, err)

	hosts, err := txeh.NewHostsDefault()
	assert.Nil(t, err)

	found, address, _ := hosts.HostAddressLookup("myhostname")
	assert.True(t, found, "'myhostname' host could not be found in /etc/hosts")
	assert.Equal(t, "127.0.0.1", address, "'myhostname' address is not equal to 127.0.0.1")
}
