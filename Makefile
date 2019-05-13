# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=urlinsane
VERSION=$(shell grep -e 'VERSION = ".*"' urlinsane.go | cut -d= -f2 | sed  s/[[:space:]]*\"//g)

.PHONY: help

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-10s\033[0m %s\n", $$1, $$2}'

all: build hash ## Build the binaries and output the md5 hash

hash: ## Output the md5 hash
	md5 builds/$(BINARY_NAME) | md5sum builds/$(BINARY_NAME)

run: ## Build and run the urlinsane tool with not options
run: deps build
	./builds/$(BINARY_NAME)

test: ## Run unit test
test: deps
	$(GOTEST) -v ./...

clean: ## Remove files created by the build
	$(GOCLEAN)
	rm -fr builds

build: ## Build the binaries for Windows, OSX, and Linux
build: deps
	mkdir -p builds
	cd cmd; $(GOBUILD) -o ../builds/$(BINARY_NAME) -v
	cd cmd; env GOOS=darwin GOARCH=amd64 $(GOBUILD) -o ../builds/$(BINARY_NAME)-$(VERSION)-darwin-amd64 -v
	cd cmd; env GOOS=linux GOARCH=amd64 $(GOBUILD) -o ../builds/$(BINARY_NAME)-$(VERSION)-linux-amd64 -v
	cd cmd; env GOOS=windows GOARCH=amd64 $(GOBUILD) -o ../builds/$(BINARY_NAME)-$(VERSION)-windows-amd64.exe -v

deps: ## Install dependencies
	$(GOGET) ./...
	$(GOGET) github.com/rangertaha/urlinsane

docker: ## Build docker image and upload to docker hub
docker: image
	docker login

image: ## Build docker image
	docker build -t urlinsane .
