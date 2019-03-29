# ahs

[![GoDoc](https://godoc.org/github.com/mvisonneau/ahs?status.svg)](https://godoc.org/github.com/mvisonneau/ahs)
[![Go Report Card](https://goreportcard.com/badge/github.com/mvisonneau/ahs)](https://goreportcard.com/report/github.com/mvisonneau/ahs)
[![Docker Pulls](https://img.shields.io/docker/pulls/mvisonneau/ahs.svg)](https://hub.docker.com/r/mvisonneau/ahs/)
[![Build Status](https://travis-ci.org/mvisonneau/ahs.svg?branch=master)](https://travis-ci.org/mvisonneau/ahs)
[![Coverage Status](https://coveralls.io/repos/github/mvisonneau/ahs/badge.svg?branch=master)](https://coveralls.io/github/mvisonneau/ahs?branch=master)

This projects aims to ease the configuration of AWS EC2 instances hostname.
In particular when they are launched as part of ASGs or fleets.

## TL;DR

```
~$ wget https://github.com/mvisonneau/ahs/releases/download/0.2.1/ahs_linux_amd64 -O /usr/local/bin/ahs; chmod +x /usr/local/bin/ahs

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

You can also use a *Dockerized version* if you prefer :

```
~$ docker run -it --rm --privileged mvisonneau/ahs <instance-id|sequential>
```

## Usage

```
~$ ahs
NAME:
   ahs - Set the hostname of an EC2 instance based on a tag value and the instance-id

USAGE:
   ahs [global options] command [command options] [arguments...]

VERSION:
   <devel>

COMMANDS:
     instance-id  compute a hostname by appending the instance-id to a prefixed/base string
     sequential   compute a sequential hostname based on the number of instances belonging to the same group
     help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --dry-run              only display what would have been done [$AHS_DRY_RUN]
   --input-tag tag        tag to use as input to determine the hostname (default: "Name") [$AHS_INPUT_TAG]
   --log-level level      log level (debug,info,warn,fatal,panic) (default: "info") [$AHS_LOG_LEVEL]
   --log-format format    log format (json,text) (default: "text") [$AHS_LOG_FORMAT]
   --output-tag tag       tag to update with the computed hostname (default: "Name") [$AHS_OUTPUT_TAG]
   --persist-hostname     set /etc/hostname with generated hostname [$AHS_PERSIST_HOSTNAME]
   --persist-hosts        assign generated hostname to 127.0.0.1 in /etc/hosts [$AHS_PERSIST_HOSTS]
   --separator separator  separator to use between tag and id (default: "-") [$AHS_SEPARATOR]
   --help, -h             show help
   --version, -v          print the version
```

### InstanceID method

This method basically takes the output of a tag considered as the `base` of the hostname, it appends a separator to it (default to `-`) and finally a truncated value of the instance-id (default to `5 characters`).

```
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

```
~$ ahs sequential -h
NAME:
   ahs sequential - compute a sequential hostname based on the number of instances belonging to the same group

USAGE:
   ahs sequential [command options]

OPTIONS:
   --instance-sequential-id-tag value  tag to which output the computed instance-sequential-id (default: "ahs:instance-id") [$AHS_INSTANCE_SEQUENTIAL_ID_TAG]
   --instance-group-tag value          tag to use in order to determine which group the instance belongs to (default: "ahs:instance-group") [$AHS_INSTANCE_GROUP_TAG]
   --respect-azs                       if instances are provisioned through an ASG, setting this flag it will get the sequential-ids associated to respective azs [$AHS_RESPECT_AZS]
```

## Develop

If you have docker locally, you can use the following command in order to quickly get a development env ready: `make dev-env`. You can also have a look onto the [Makefile](/Makefile) in order to see all available options:

```
~$ make
all                            Test, builds and ship package for all supported platforms
build                          Build the binary
clean                          Remove binary if it exists
coverage                       Generates coverage report
dev-env                        Build a local development environment using Docker
fmt                            Format source code
help                           Displays this help
imports                        Fixes the syntax (linting) of the codebase
install                        Build and install locally the binary (dev purpose)
lint                           Run golint and go vet against the codebase
publish-github                 Send the binaries onto the GitHub release
setup                          Install required libraries/tools
test                           Run the tests against the codebase
```

## Contribute

Contributions are more than welcome! Feel free to submit a [PR](https://github.com/mvisonneau/ahs/pulls).
