# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed

- Upgrade the build toolchain to Go 1.27.4 and update Go dependencies.
- Use an Alpine-based Go image for Docker builds. By default, `make build` targets the Docker image's architecture.

### Fixed

- Prevent a divide-by-zero panic when `-time` specifies a positive duration shorter than one millisecond.
- Start the initial `-workers` processes concurrently so they can ramp up memory together.

[unreleased]: https://github.com/chaos-mesh/memStress/compare/v0.3.1...HEAD
