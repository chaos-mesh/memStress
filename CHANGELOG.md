# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Removed

- Remove percentage values from `-size`. Use a byte size such as `64MiB`. Chaos Daemon converts percentage-based StressChaos settings to bytes.

## [0.4.0] - 2026-11-19

### Changed

- Upgrade the build toolchain to Go 1.27.4 and update Go dependencies. ([#13](https://github.com/chaos-mesh/memStress/pull/13))
- Use an Alpine-based Go image for Docker builds. By default, `make build` targets the Docker image's architecture. ([#15](https://github.com/chaos-mesh/memStress/pull/15))

### Fixed

- Prevent a divide-by-zero panic when `-time` specifies a positive duration shorter than one millisecond. ([#11](https://github.com/chaos-mesh/memStress/pull/11))
- Start the initial `-workers` processes concurrently so they can ramp up memory together. ([#12](https://github.com/chaos-mesh/memStress/pull/12))

[unreleased]: https://github.com/chaos-mesh/memStress/compare/v0.4.0...HEAD
[0.4.0]: https://github.com/chaos-mesh/memStress/compare/v0.3.1...v0.4.0
