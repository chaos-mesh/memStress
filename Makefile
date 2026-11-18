# Build, vet, and release use the Go Alpine image by default.
# Keep GO_VERSION in sync with go.mod.
GO_VERSION ?= 1.27.4
VERSION    ?= dev
DIST       ?= dist

RUN = docker run --rm -e CGO_ENABLED=0 \
	-v "$(CURDIR)":/src -w /src golang:$(GO_VERSION)-alpine

.DEFAULT_GOAL := help
.PHONY: help build vet release clean

help:
	@echo "make build                   build ./memStress for the Docker image's architecture"
	@echo "make release VERSION=v0.3.1  package static Linux amd64 and arm64 binaries in $(DIST)/"
	@echo "make vet                     run go vet ./..."
	@echo "make clean                   remove ./memStress and $(DIST)/"

build:
	$(RUN) go build -trimpath -ldflags="-s -w" -o memStress main.go

vet:
	$(RUN) go vet ./...

# Keep the archive names and root-level memStress binary compatible with chaos-daemon.
release:
	$(RUN) sh -euc 'dist="$$1"; version="$$2"; \
	  mkdir -p "$$dist"; \
	  build_dir=$$(mktemp -d); \
	  trap "rm -rf \"$$build_dir\"" EXIT; \
	  package() { \
	    GOARCH="$$1" go build -trimpath -ldflags="-s -w" -o "$$build_dir/memStress" main.go; \
	    tar -czf "$$dist/memStress_$$version-$$2-linux-gnu.tar.gz" -C "$$build_dir" memStress; \
	  }; \
	  package amd64 x86_64; \
	  package arm64 aarch64' sh "$(DIST)" "$(VERSION)"
	@ls -l "$(DIST)"

clean:
	rm -rf "$(DIST)" memStress
