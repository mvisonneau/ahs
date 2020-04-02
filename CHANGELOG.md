# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [0ver](https://0ver.org).

## [Unreleased]

### Added

- Release binaries using goreleaser
- Releases 386 and arm64 arch
- deb and rpm packages

### Changed

- Bumped to go 1.14
- Bumped all go modules to their latest versions

## [0.2.3] - 2019-05-28

### Added

- Ensure that we do not have sequential id higher than the ASG max size

### Changed

- Fix fmt test function false positive
- Bumped dependencies

## [0.2.2] - 2019-03-30

### Added

- Release binaries are now automatically built and published from the CI

### Changed

- Moved CI from `Travis` to `Drone`
- Optimized Makefile

## [0.2.1] - 2019-03-29

### Added

- New flag `--persist-hostname` which will update the `/etc/hostname` file value with the generated hostname
- New flag `--persist-hosts` which will set the generated hostname as a host pointing to the loopback
- Released `arm64` binaries

### Changed

- Makefile improvements
- Updated dependencies to their latest versions
- Updated to `go 1.12`
- Switched to **go modules**
- Updated `Dockerfile` to use **busybox** instead of **scratch** image as source

## [0.2.0] - 2018-07-25

### Added

- Added a new flag `--respect-azs` for `sequential` method that ensure we keep sequential ids aligned with configured ASG AZs.

### Changed

- Some coverage tweaks/cleanup

## [0.1.1] - 2018-07-23

### Changed

- Fixed the `dry-run` function
- Filter on running instances only for the **sequential** method
- Avoid duplicates when looking for sequential ids which breaks the compute function

## [0.1.0] - 2018-07-23

### Added

- New **sequential** hostname calculation method

### Changed

- Ensure that we are running as **root**
- Refactored the codebase and added more parameters
- Updated all deps to most recent version

## [0.0.3] - 2018-07-13

### Changed

- Boot issues when the tag is not available through the API yet [GH-1]
- Fixed the Dockerized version, also updated doc
- Updated dependencies to latest version and removed some constraints
- Disabled CGO on build function

## [0.0.2] - 2018-06-18

### Added

- Added some tests
- Added CI config

### Changed

- Do not keep appending the instance-id when it is already set on the inputTag
- Fixed incorrect log output on output tag value

## [0.0.1] - 2018-06-13

### Added

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

[Unreleased]: https://github.com/mvisonneau/ahs/compare/0.2.3...HEAD
[0.2.3]: https://github.com/mvisonneau/ahs/tree/0.2.3
[0.2.2]: https://github.com/mvisonneau/ahs/tree/0.2.2
[0.2.1]: https://github.com/mvisonneau/ahs/tree/0.2.1
[0.2.0]: https://github.com/mvisonneau/ahs/tree/0.2.0
[0.1.1]: https://github.com/mvisonneau/ahs/tree/0.1.1
[0.1.0]: https://github.com/mvisonneau/ahs/tree/0.1.0
[0.0.3]: https://github.com/mvisonneau/ahs/tree/0.0.3
[0.0.2]: https://github.com/mvisonneau/ahs/tree/0.0.2
[0.0.1]: https://github.com/mvisonneau/ahs/tree/0.0.1
