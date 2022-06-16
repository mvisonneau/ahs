# 🏷 ahs - AWS EC2 Hostname Setter

[![PkgGoDev](https://pkg.go.dev/badge/github.com/mvisonneau/ahs)](https://pkg.go.dev/mod/github.com/mvisonneau/ahs)
[![Go Report Card](https://goreportcard.com/badge/github.com/mvisonneau/ahs)](https://goreportcard.com/report/github.com/mvisonneau/ahs)
[![Docker Pulls](https://img.shields.io/docker/pulls/mvisonneau/ahs.svg)](https://hub.docker.com/r/mvisonneau/ahs/)
[![Build Status](https://cloud.drone.io/api/badges/mvisonneau/ahs/status.svg)](https://cloud.drone.io/mvisonneau/ahs)
[![Coverage Status](https://coveralls.io/repos/github/mvisonneau/ahs/badge.svg?branch=master)](https://coveralls.io/github/mvisonneau/ahs?branch=master)

This projects aims to ease the configuration of AWS EC2 instances hostname.
In particular when they are launched as part of ASGs or fleets.

## TL;DR

```bash
# Using instance-id method
~$ ahs instance-id
INFO[2018-07-23T11:56:00Z] Found AZ: 'eu-west-1a'
INFO[2018-07-23T11:56:00Z] Computed region : 'eu-west-1'
INFO[2018-07-23T11:56:00Z] Found instance-id : 'i-096bed3161783f000'
INFO[2018-07-23T11:56:00Z] Querying input-tag 'Name' from EC2 API
INFO[2018-07-23T11:56:00Z] Computing hostname with truncated instance-id
INFO[2018-07-23T11:56:00Z] Computed unique hostname : 'myhostname-096be'
INFO[2018-07-23T11:56:00Z] Setting instance hostname locally
INFO[2018-07-23T11:56:00Z] Setting hostname on configured instance output tag 'Name'

# Using sequential method
~$ ahs sequential
INFO[2018-07-23T11:56:00Z] Found AZ: 'eu-west-1a'
INFO[2018-07-23T11:56:00Z] Computed region : 'eu-west-1'
INFO[2018-07-23T11:56:00Z] Found instance-id : 'i-096bed3161783f000'
INFO[2018-07-23T11:56:00Z] Querying input-tag 'Name' from EC2 API
INFO[2018-07-23T11:56:00Z] Computing a hostname with sequential naming
INFO[2018-07-23T11:56:00Z] Computed unique hostname : 'myhostname-1' - Sequential ID : '1'
INFO[2018-07-23T11:56:00Z] Setting instance hostname locally
INFO[2018-07-23T11:56:00Z] Setting hostname on configured instance output tag 'Name'
INFO[2018-07-23T11:56:00Z] Setting instance sequential id (1) on configured tag 'ahs:instance-id'
```

You can also use a containerized version if you prefer:

```bash
~$ docker run -it --rm --privileged mvisonneau/ahs <instance-id|sequential>
```

## Install

Have a look onto the [latest release page](https://github.com/mvisonneau/ahs/releases/latest) and pick your flavor.

### Go

```bash
~$ go install github.com/mvisonneau/ahs/cmd/ahs@latest
```

### Homebrew

```bash
~$ brew install mvisonneau/tap/ahs
```

### Docker

```bash
~$ docker run -it --rm --privileged mvisonneau/ahs
```

### Binaries, DEB and RPM packages

For the following ones, you need to know which version you want to install, to fetch the latest available :

```bash
~$ export AHS_VERSION=$(curl -s "https://api.github.com/repos/mvisonneau/ahs/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
```

```bash
# Binary (eg: linux/amd64)
~$ wget https://github.com/mvisonneau/ahs/releases/download/${AHS_VERSION}/ahs_${AHS_VERSION}_linux_amd64.tar.gz
~$ tar zxvf ahs_${AHS_VERSION}_linux_amd64.tar.gz -C /usr/local/bin

# DEB package (eg: linux/386)
~$ wget https://github.com/mvisonneau/ahs/releases/download/${AHS_VERSION}/ahs_${AHS_VERSION}_linux_386.deb
~$ dpkg -i ahs_${AHS_VERSION}_linux_386.deb

# RPM package (eg: linux/arm64)
~$ wget https://github.com/mvisonneau/ahs/releases/download/${AHS_VERSION}/ahs_${AHS_VERSION}_linux_arm64.rpm
~$ rpm -ivh ahs_${AHS_VERSION}_linux_arm64.rpm
```

## Usage

```bash
~$ ahs
NAME:
   ahs - Set the hostname of an EC2 instance based on a tag value and the instance-id

USAGE:
   ahs [global options] command [command options] [arguments...]

COMMANDS:
   instance-id  compute a hostname by appending the instance-id to a prefixed/base string
   sequential   compute a sequential hostname based on the number of instances belonging to the same group
   help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --dry-run              only display what would have been done (default: false) [$AHS_DRY_RUN]
   --help, -h             show help (default: false)
   --input-tag tag        tag to use as input to determine the hostname (default: "Name") [$AHS_INPUT_TAG]
   --log-format format    log format (json,text) (default: "text") [$AHS_LOG_FORMAT]
   --log-level level      log level (debug,info,warn,fatal,panic) (default: "info") [$AHS_LOG_LEVEL]
   --output-tag tag       tag to update with the computed hostname (default: "Name") [$AHS_OUTPUT_TAG]
   --persist-hostname     set /etc/hostname with generated hostname (default: false) [$AHS_PERSIST_HOSTNAME]
   --persist-hosts        assign generated hostname to 127.0.0.1 in /etc/hosts (default: false) [$AHS_PERSIST_HOSTS]
   --separator separator  separator to use between tag and id (default: "-") [$AHS_SEPARATOR]
   --version, -v          print the version
```

### InstanceID method

This method basically takes the output of a tag considered as the `base` of the hostname, it appends a separator to it (default to `-`) and finally a truncated value of the instance-id (default to `5 characters`).

```bash
~$ ahs instance-id -h
NAME:
   ahs instance-id - compute a hostname by appending the instance-id to a prefixed/base string

USAGE:
   ahs instance-id [command options]  

OPTIONS:
   --length value  length of the id to keep in the hostname (default: 5) [$AHS_INSTANCE_ID_LENGTH]
```

### Sequential method

This method allows you to have sequential hostnames on instances on which you couldn't or haven't configured at the time of provisioning.

```bash
~$ ahs sequential -h
NAME:
   ahs sequential - compute a sequential hostname based on the number of instances belonging to the same group

USAGE:
   ahs sequential [command options]  

OPTIONS:
   --instance-group-tag value          tag to use in order to determine which group the instance belongs to (default: "ahs:instance-group") [$AHS_INSTANCE_GROUP_TAG]
   --instance-sequential-id-tag value  tag to which output the computed instance-sequential-id (default: "ahs:instance-id") [$AHS_INSTANCE_SEQUENTIAL_ID_TAG]
   --respect-azs                       if instances are provisioned through an ASG, setting this flag it will get the sequential-ids associated to respective azs (default: false) [$AHS_RESPECT_AZS]
```

## Develop

You can have a look at the [Makefile](/Makefile) in order to see all available options:

```bash
~$ make
all                            Test, builds and ship package for all supported platforms
build                          Build the binaries using local GOOS
clean                          Remove binary if it exists
coverage                       Generates coverage report
coverage-html                  Generates coverage report and displays it in the browser
fmt                            Format source code
help                           Displays this help
install                        Build and install locally the binary (dev purpose)
is-git-dirty                   Tests if git is in a dirty state
lint                           Run all lint related tests upon the codebase
prerelease                     Build & prerelease the binaries (edge)
release                        Build & release the binaries (stable)
setup                          Install required libraries/tools for build tasks
test                           Run the tests against the codebase
```

## Contribute

Contributions are more than welcome! Feel free to submit a [PR](https://github.com/mvisonneau/ahs/pulls).
