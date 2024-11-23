# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GODOC=$(GOCMD)doc
BDIR=build
BINARY_NAME=urlinsane
VERSION=$(shell grep -e 'VERSION = ".*"' internal/version.go | cut -d= -f2 | sed  s/[[:space:]]*\"//g)

.PHONY: help version build dpkg deps test clean doc



help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-10s\033[0m %s\n", $$1, $$2}'

version: ## Returns the version number
	@echo $(VERSION)


build: deps ## Build the binaries for Windows, OSX, and Linux
	mkdir -p $(BDIR)
	# $(GOBUILD) -C cmd -o ../$(BDIR)/$(BINARY_NAME) -v
	# env GOOS=darwin GOARCH=amd64 $(GOBUILD) -C cmd -o ../$(BDIR)/$(BINARY_NAME)-$(VERSION)-darwin-amd64 -v
	# sha512sum $(BDIR)/$(BINARY_NAME)-$(VERSION)-darwin-amd64 > $(BDIR)/$(BINARY_NAME)-$(VERSION)-darwin-amd64.sha512

	# env GOOS=linux GOARCH=amd64 $(GOBUILD) -C cmd -o ../$(BDIR)/$(BINARY_NAME)-$(VERSION)-linux-amd64 -v
	# sha512sum $(BDIR)/$(BINARY_NAME)-$(VERSION)-linux-amd64 > $(BDIR)/$(BINARY_NAME)-$(VERSION)-linux-amd64.sha512

	# env GOOS=windows GOARCH=amd64 $(GOBUILD) -C cmd -o ../$(BDIR)/$(BINARY_NAME)-$(VERSION)-windows-amd64.exe -v
	# sha512sum $(BDIR)/$(BINARY_NAME)-$(VERSION)-windows-amd64.exe > $(BDIR)/$(BINARY_NAME)-$(VERSION)-windows-amd64.exe.sha512


install: build ## Install the binaries in Linux
	mkdir -p $(BDIR)


deps: ## Install dependencies
	$(GOGET) ./...

test: deps ## Run unit test
	$(GOTEST) -v ./...

clean: ## Remove files created by the build
	$(GOCLEAN)
	rm -rf build
	

doc: ## Go documentation
	$(GODOC) -http=:6060


update: ## Update data files
	bash scripts/update.sh
