# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed

- Release binaries now include darwin/amd64, darwin/arm64, windows/amd64, and windows/arm64 alongside the existing linux targets. Windows binaries are named `waluigi-windows-<arch>.exe`.

## [0.2.0] - 2025-06-22

### Added

- Add support for `aws-resolver-rules-operator` logs.

## [0.1.0] - 2025-05-12

### Added

- Add new flag `level` to filter which logs to show.

[Unreleased]: https://github.com/giantswarm/waluigi/compare/v0.2.0...HEAD
[0.2.0]: https://github.com/giantswarm/waluigi/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/giantswarm/waluigi/compare/v0.0.1...v0.1.0
[0.0.1]: https://github.com/giantswarm/cluster-aws/releases/tag/v0.0.1
