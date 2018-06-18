# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]

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

[Unreleased]: https://github.com/mvisonneau/ahs/compare/0.0.2...HEAD
[0.0.2]: https://github.com/mvisonneau/ahs/tree/0.0.2
[0.0.1]: https://github.com/mvisonneau/ahs/tree/0.0.1
