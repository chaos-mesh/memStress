# Repository Guidelines

## Project Structure & Module Organization

memStress is a Linux command-line tool that simulates memory allocation. The root package contains `main.go` for flag parsing, worker subprocesses, and memory allocation, and `main_test.go` for unit tests. The `Makefile` defines build and packaging targets. GitHub Actions workflows live in `.github/workflows/`.

## Build, Test, and Development Commands

- `make help` lists available targets.
- `make build` builds a static Linux amd64 executable at `./memStress` using Docker. Keep the Makefile's `GO_VERSION` aligned with `go.mod`.
- `make vet` runs `go vet ./...` in the same container.
- `make release VERSION=vX.Y.Z` packages Linux amd64 and arm64 binaries into separate tarballs under `dist/`. CI runs vet and builds both packages.
- On Linux, `PATH="$PWD:$PATH" ./memStress -size 16MiB -time 1s -workers 2` ramps up to 16 MiB total across two workers. The parent finds worker binaries on `PATH`. Press Ctrl+C to stop it after the ramp.

Reuse existing Go build and module caches and existing Docker image and build caches when available. Avoid creating per-task caches.

## Coding Style & Naming Conventions

Format Go changes with `gofmt` and use its standard indentation. Keep standard-library and third-party imports in separate groups. Follow existing Go naming conventions, such as `linearGrow` for an unexported function, and retain Apache 2.0 headers in Go source files.

## Testing Guidelines

Use Go's standard `testing` package. Put tests in `*_test.go` files and name them `Test<Behavior>`, following `TestLinearGrowSubMillisecond`. Run `go test ./...` on Linux. Add targeted regression tests for allocation or worker changes. CI currently runs vet and packaging checks, with no test coverage threshold configured.

## Commit & Pull Request Guidelines

Follow recent commit subjects with a short imperative summary and a prefix such as `fix:`, `fix(workers):`, or `chore:`. In pull requests, explain the behavior change, link related issues when available, and report the checks you ran.
