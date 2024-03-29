package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type inputs struct {
	hostname   string
	separator  string
	instanceID string
	length     int
}

func TestComputeHostnameWithInstanceID(t *testing.T) {
	tests := map[string]struct {
		expectedResult string

		inputs inputs
	}{
		"truncated hostname":                                 {"myhostname-12345", inputs{"myhostname", "-", "i-123456789", 5}},
		"truncated hostname second-run":                      {"myhostname-12345", inputs{"myhostname-12345", "-", "i-123456789", 5}},
		"truncated hostname expanded (length > instance id)": {"myhostname-123456789", inputs{"myhostname", "-", "i-123456789", 100}},

		"hostname expanded (full-length)": {"myhostname-123456789", inputs{"myhostname", "-", "i-123456789", -1}},
		"kebab hostname second-run":       {"my-host-name-12345", inputs{"my-host-name-12345", "-", "i-123456789", 5}},
		"kebab hostname expanded":         {"my-host-name-12345", inputs{"my-host-name", "-", "i-123456789", 5}},

		"kebab hostname (full-length)":           {"my-host-name-123456789", inputs{"my-host-name", "-", "i-123456789", -1}},
		"kebab hostname expanded (full-length)":  {"my-host-name-12-123456789", inputs{"my-host-name-12", "-", "i-123456789", -1}},
		"kebab duplicate id hostname truncated7": {"my-abcdefg-host-abcdefg", inputs{"my-abcdefg-host", "-", "i-abcdefg", -1}},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			hostname, err := computeHostnameWithInstanceID(tt.inputs.hostname,
				tt.inputs.separator,
				tt.inputs.instanceID,
				tt.inputs.length)
			assert.Nil(t, err)
			assert.Equal(t, tt.expectedResult, hostname)
		})
	}
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

// func TestUpdateHostnameFile(t *testing.T) {
// 	err := updateHostnameFile("myhostname")
// 	assert.Nil(t, err)

// 	content, err := ioutil.ReadFile("/etc/hostname")
// 	assert.Nil(t, err)
// 	assert.Equal(t, "myhostname\n", string(content))
// }

// func TestUpdateHostsFile(t *testing.T) {
// 	err := updateHostsFile("myhostname")
// 	assert.Nil(t, err)

// 	hosts, err := txeh.NewHostsDefault()
// 	assert.Nil(t, err)

// 	found, address, _ := hosts.HostAddressLookup("myhostname")
// 	assert.True(t, found, "'myhostname' host could not be found in /etc/hosts")
// 	assert.Equal(t, "127.0.0.1", address, "'myhostname' address is not equal to 127.0.0.1")
// }
