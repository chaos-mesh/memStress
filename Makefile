# memStress is Linux-only, so every target runs inside a pinned Go container.
# This builds on any host (incl. macOS) with just Docker, always uses the same
# Go toolchain, and produces a static linux/amd64 binary. Keep GO_VERSION in
# sync with go.mod.
GO_VERSION ?= 1.25.11
VERSION    ?= dev
DIST       ?= dist

RUN = docker run --rm --platform linux/amd64 -e CGO_ENABLED=0 \
	-v "$(CURDIR)":/src -w /src golang:$(GO_VERSION)

.DEFAULT_GOAL := help
.PHONY: help build vet release clean

help:
	@echo "make build                   build ./memStress (linux/amd64)"
	@echo "make release VERSION=v0.3.1  build release tarballs into ./$(DIST)"
	@echo "make vet                     run go vet"
	@echo "make clean                   remove build artifacts"

build:
	$(RUN) go build -trimpath -ldflags="-s -w" -o memStress main.go

vet:
	$(RUN) go vet ./...

# Cross-compile both arches and package each as a tarball containing a single
# `memStress` binary (named to match what chaos-mesh's chaos-daemon expects).
release:
	$(RUN) sh -euc 'mkdir -p $(DIST); \
	  GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o $(DIST)/memStress main.go; \
	  tar -czf "$(DIST)/memStress_$(VERSION)-x86_64-linux-gnu.tar.gz"  -C $(DIST) memStress; \
	  GOARCH=arm64 go build -trimpath -ldflags="-s -w" -o $(DIST)/memStress main.go; \
	  tar -czf "$(DIST)/memStress_$(VERSION)-aarch64-linux-gnu.tar.gz" -C $(DIST) memStress; \
	  rm -f $(DIST)/memStress'
	@ls -l $(DIST)

clean:
	rm -rf $(DIST) memStress
