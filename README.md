# ahs

This projects aims to ease the configuration of AWS EC2 instances hostname.
In particular when they are launched as part of ASGs for fleets.

## Usage

```
~$ ahs
NAME:
   ahs - Set the hostname of an EC2 instance based on a tag value and the instance-id

USAGE:
   ahs [global options] command [command options] [arguments...]

COMMANDS:
     run      replace the hostname with found/computed values
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --dry-run           only display what would have been done [$AHS_DRY_RUN]
   --id-length value   length of the id to keep in the hostname (default: 5) [$AHS_ID_LENGTH]
   --input-tag value   tag to use as input to determine the hostname (default: "Name") [$AHS_TAG_NAME_INPUT]
   --log-level value   log level (debug,info,warn,fatal,panic) (default: "info") [$AHS_LOG_LEVEL]
   --log-format value  log format (json,text) (default: "text") [$AHS_LOG_FORMAT]
   --output-tag value  tag to update with the computed hostname (default: "Name") [$AHS_TAG_NAME_OUTPUT]
   --separator value   separator to use between tag and id (default: "-") [$AHS_SEPARATOR]
   --help, -h          show help
   --version, -v       print the version
```
