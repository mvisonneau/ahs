# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]
### BUGFIXES
- Fixed the `dry-run` function
- Filter on running instances only for the **sequential** method

## [0.1.0] - 2018-07-23
### FEATURES
- New **sequential** hostname calculation method
- Ensure that we are running as **root**

### ENHANCEMENTS
- Refactored the codebase and added more parameters
- Updated all deps to most recent version

## [0.0.3] - 2018-07-13
### BUGFIXES
- Boot issues when the tag is not available through the API yet [GH-1]
- Fixed the Dockerized version, also updated doc

### ENHANCEMENTS
- Updated dependencies to latest version and removed some constraints
- Disabled CGO on build function

## [0.0.2] - 2018-06-18
### FEATURES
- Do not keep appending the instance-id when it is already set on the inputTag

### ENHANCEMENTS
- Added some tests
- Added CI config

### BUGFIXES
- Fixed incorrect log output on output tag value

## [0.0.1] - 2018-06-13
### FEATURES
- Working state of the app
- Configure hostnames on Unix based OSes
- Hostname based on a input tag and the instance-id
- Configurable separator
- Configurable length of the instance-id to include in the hostname
- Configurable input and output tags
- dry-run capability
- Makefile
- License
- Readme

[Unreleased]: https://github.com/mvisonneau/ahs/compare/0.1.0...HEAD
[0.1.0]: https://github.com/mvisonneau/ahs/tree/0.1.0
[0.0.3]: https://github.com/mvisonneau/ahs/tree/0.0.3
[0.0.2]: https://github.com/mvisonneau/ahs/tree/0.0.2
[0.0.1]: https://github.com/mvisonneau/ahs/tree/0.0.1
