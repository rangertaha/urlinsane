# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GODOC=$(GOCMD) doc
BDIR=build
$(shell mkdir -p $(BDIR))
BINARY_NAME=urlinsane
VERSION=$(shell grep -e 'VERSION = ".*"' internal/version.go | cut -d= -f2 | sed  s/[[:space:]]*\"//g)
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)

# -trimpath removes local filesystem paths from the binary, so the same source
# builds byte-identically on another machine. -s -w drop the symbol table and
# DWARF; the data this embeds is the size, not the debug info, but there is no
# reason to ship both.
LDFLAGS=-s -w -X github.com/rangertaha/urlinsane/internal.COMMIT=$(COMMIT)
BUILDFLAGS=-trimpath -ldflags "$(LDFLAGS)"

.PHONY: help version build install dpkg deps test race vet fmt check dataset clean doc release dist update

# Every package, so a new top-level directory is covered without editing this.
PKGS=./...



help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-10s\033[0m %s\n", $$1, $$2}'

version: ## Returns the version number
	@echo $(VERSION)


# Every target Go supports that this tool makes sense on. Cross-compiling is
# only possible because the sqlite driver is pure Go: with the cgo driver each
# of these needed its own C toolchain, and CGO_ENABLED=0 built a binary that
# ran but could not open the database.
PLATFORMS=\
	darwin/amd64 darwin/arm64 \
	linux/amd64 linux/arm64 linux/386 linux/arm linux/riscv64 \
	windows/amd64 windows/arm64 windows/386 \
	freebsd/amd64 freebsd/arm64 \
	openbsd/amd64 netbsd/amd64

release: deps ## Build release binaries for every platform
	@set -e; for p in $(PLATFORMS); do \
		os=$${p%/*}; arch=$${p#*/}; \
		ext=""; [ "$$os" = "windows" ] && ext=".exe"; \
		out="$(BDIR)/$(BINARY_NAME)-$(VERSION)-$$os-$$arch$$ext"; \
		echo "  $$os/$$arch"; \
		CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch \
			$(GOBUILD) $(BUILDFLAGS) -o "$$out" ./cmd/$(BINARY_NAME); \
		sha512sum "$$out" > "$$out.sha512"; \
	done
	@echo "\n$(BDIR):"; ls -1 $(BDIR) | grep -v sha512

# Archives rather than bare binaries. A release asset that unpacks to a
# directory carrying the licence is what a human and a package manager both
# expect, and one checksums.txt over the archives is easier to verify than a
# .sha512 beside every binary.
#
# Depends on release rather than repeating the build, so the archives can only
# ever contain the binaries release just produced.
dist: release ## Package the release binaries into per-platform archives
	@set -e; rm -rf $(BDIR)/dist; mkdir -p $(BDIR)/dist; \
	for p in $(PLATFORMS); do \
		os=$${p%/*}; arch=$${p#*/}; \
		ext=""; [ "$$os" = "windows" ] && ext=".exe"; \
		name="$(BINARY_NAME)-$(VERSION)-$$os-$$arch"; \
		stage="$(BDIR)/dist/$$name"; \
		mkdir -p "$$stage"; \
		cp "$(BDIR)/$$name$$ext" "$$stage/$(BINARY_NAME)$$ext"; \
		cp LICENSE README.md CHANGELOG.md "$$stage/"; \
		if [ "$$os" = "windows" ]; then \
			( cd $(BDIR)/dist && zip -qr "$$name.zip" "$$name" ); \
		else \
			tar -czf "$(BDIR)/dist/$$name.tar.gz" -C "$(BDIR)/dist" "$$name"; \
		fi; \
		rm -rf "$$stage"; \
	done; \
	( cd $(BDIR)/dist && sha256sum *.tar.gz *.zip > checksums.txt )
	@echo; ls -1 $(BDIR)/dist

build: deps ## Build both binaries
	$(GOBUILD) -C cmd/$(BINARY_NAME) -o ../../$(BDIR)/$(BINARY_NAME)
	$(GOBUILD) -o $(BDIR)/datasets ./cmd/datasets

install: build ## Install the binaries in Linux
	@chmod +x $(BDIR)/$(BINARY_NAME)
	@sudo mv $(BDIR)/$(BINARY_NAME) /usr/local/bin/

deps: ## Download module dependencies
	# go mod download, not `go get ./...`: in module mode go get mutates go.mod,
	# so a build target that ran it could silently change the dependency graph.
	$(GOCMD) mod download

test: deps ## Run unit tests
	$(GOTEST) $(PKGS)

race: deps ## Run tests under the race detector
	# The scheduler runs operators concurrently and merges their deltas through
	# one applier, so a data race here is a wrong graph rather than a crash.
	# This is the run that would catch it.
	$(GOTEST) -race $(PKGS)

vet: ## Report suspicious constructs
	$(GOCMD) vet $(PKGS)

fmt: ## Report unformatted files, without rewriting them
	@out=$$(gofmt -l cmd internal pkg datasets); \
	if [ -n "$$out" ]; then echo "unformatted:"; echo "$$out"; exit 1; fi

check: fmt vet race ## Everything CI should run

dataset: ## Rebuild the embedded reference database from datasets/
	# The output is internal/config/dataset.db, which //go:embed picks up, so
	# this must run before build for a data change to reach the binary.
	$(GOCMD) run ./cmd/datasets build datasets

clean: ## Remove build output
	$(GOCLEAN)
	rm -rf $(BDIR)
	
doc: ## Go documentation
	$(GODOC) -http=:6060

update: ## Update data files
	bash scripts/update.sh

dpkg:  ## Build debian package
	# dpkg-buildpackage -b -rfakeroot -us -uc
	debuild  -us -uc
