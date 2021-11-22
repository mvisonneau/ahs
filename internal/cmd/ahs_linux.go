//go:build linux
// +build linux

package cmd

import "syscall"

func setSystemHostname(hostname string) error {
	return syscall.Sethostname([]byte(hostname))
}
