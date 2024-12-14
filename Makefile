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

.PHONY: help version build install dpkg deps test clean doc



help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-10s\033[0m %s\n", $$1, $$2}'

version: ## Returns the version number
	@echo $(VERSION)


release: deps ## Build the binaries for Windows, OSX, and Linux
	env GOOS=darwin GOARCH=amd64 $(GOBUILD) -C cmd/$(BINARY_NAME) -o ../../$(BDIR)/$(BINARY_NAME)-$(VERSION)-darwin-amd64 -v
	sha512sum $(BDIR)/$(BINARY_NAME)-$(VERSION)-darwin-amd64 > $(BDIR)/$(BINARY_NAME)-$(VERSION)-darwin-amd64.sha512

	env GOOS=linux GOARCH=amd64 $(GOBUILD) -C cmd/$(BINARY_NAME) -o ../../$(BDIR)/$(BINARY_NAME)-$(VERSION)-linux-amd64 -v
	sha512sum $(BDIR)/$(BINARY_NAME)-$(VERSION)-linux-amd64 > $(BDIR)/$(BINARY_NAME)-$(VERSION)-linux-amd64.sha512

	env GOOS=windows GOARCH=amd64 $(GOBUILD) -C cmd/$(BINARY_NAME) -o ../../$(BDIR)/$(BINARY_NAME)-$(VERSION)-windows-amd64.exe -v
	sha512sum $(BDIR)/$(BINARY_NAME)-$(VERSION)-windows-amd64.exe > $(BDIR)/$(BINARY_NAME)-$(VERSION)-windows-amd64.exe.sha512

build: deps ## Build the binary
	$(GOBUILD) -C cmd/$(BINARY_NAME) -o ../../$(BDIR)/$(BINARY_NAME)

# install: deps ## Install the binaries in Linux
# @mkdir -p $(BDIR)
# $(GOBUILD) -C cmd -o ../$(BDIR)/$(BINARY_NAME)
# @chmod +x $(BDIR)/$(BINARY_NAME)
# @sudo mv $(BDIR)/$(BINARY_NAME) /usr/local/bin/


deps: ## Install dependencies
	$(GOGET) ./...

# test: deps ## Run unit test
# 	# $(GOTEST) -v ./internal/... ./cmd/...
# 	go test -v ./internal/... ./cmd/...

clean: ## Remove files build files
	$(GOCLEAN)
	rm -rf build
	
doc: ## Go documentation
	$(GODOC) -http=:6060

update: ## Update data files
	bash scripts/update.sh

dpkg:  ## Build debian package
	# dpkg-buildpackage -b -rfakeroot -us -uc
	debuild  -us -uc
