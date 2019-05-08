# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=urlinsane
VERSION=0.5.1

.PHONY: help

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-10s\033[0m %s\n", $$1, $$2}'


all: build hash ## Build a binary and output the md5 hash

hash: ## Output the md5 hash
	md5 builds/$(BINARY_NAME) | md5sum builds/$(BINARY_NAME)

build: ## Build the binary in /builds
	mkdir -p builds
	cd cmd; $(GOBUILD) -o ../builds/$(BINARY_NAME) -v

test: ## Run unit test
	$(GOTEST) -v ./...

clean: ## Remove files created by the build
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

run: ## Build and run the urlinsane tool
	cd cmd; $(GOBUILD) -o ../$(BINARY_NAME) -v
	./$(BINARY_NAME)

deps: ## Install dependencies
	$(GOGET) github.com/rangertaha/urlinsane

versions: ## Build the binaries for Windows, OSX, and Linux
	mkdir -p builds
	cd cmd; env GOOS=darwin GOARCH=amd64 $(GOBUILD) -o ../builds/$(BINARY_NAME)-$(VERSION)-darwin-amd64 -v
	cd cmd; env GOOS=linux GOARCH=amd64 $(GOBUILD) -o ../builds/$(BINARY_NAME)-$(VERSION)-linux-amd64 -v
	cd cmd; env GOOS=windows GOARCH=amd64 $(GOBUILD) -o ../builds/$(BINARY_NAME)-$(VERSION)-windows-amd64.exe -v
