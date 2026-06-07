GO ?= go
LDFLAGS ?= -s -w
VERSION ?= dev
DIST ?= dist

# <GOARCH>:<release arch name> pairs. The release arch name matches the suffix
# used by the release tarballs (and expected by chaos-mesh's chaos-daemon image).
ARCHS := amd64:x86_64 arm64:aarch64

.PHONY: build
build:
	$(GO) build -trimpath -ldflags="$(LDFLAGS)" -o memStress main.go

.PHONY: vet
vet:
	GOOS=linux $(GO) vet ./...

# Cross-compile static linux binaries for every supported arch and package them
# as dist/memStress_$(VERSION)-<arch>-linux-gnu.tar.gz (each tarball contains a
# single `memStress` binary at its root).
.PHONY: release
release:
	rm -rf $(DIST)
	mkdir -p $(DIST)
	@for pair in $(ARCHS); do \
		goarch=$${pair%%:*}; arch=$${pair##*:}; \
		echo "==> building memStress $(VERSION) for linux/$$arch ($$goarch)"; \
		CGO_ENABLED=0 GOOS=linux GOARCH=$$goarch $(GO) build -trimpath -ldflags="$(LDFLAGS)" -o $(DIST)/memStress main.go; \
		tar -czf $(DIST)/memStress_$(VERSION)-$$arch-linux-gnu.tar.gz -C $(DIST) memStress; \
		rm -f $(DIST)/memStress; \
	done
	@ls -l $(DIST)

.PHONY: clean
clean:
	rm -rf $(DIST) memStress
