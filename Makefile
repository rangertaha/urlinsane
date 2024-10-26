# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GODOC=$(GOCMD)doc
BINARY_NAME=urlinsane
# GCP_PROJECT_ID=cyberse
# GCR_HOST=gcr.io
VERSION=$(shell grep -e 'VERSION = ".*"' version.go | cut -d= -f2 | sed  s/[[:space:]]*\"//g)

.PHONY: help

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-10s\033[0m %s\n", $$1, $$2}'

build: deps ## Build the binaries for Windows, OSX, and Linux
	mkdir -p builds
	cd cmd; $(GOBUILD) -o ../builds/$(BINARY_NAME) -v
	cd cmd; env GOOS=darwin GOARCH=amd64 $(GOBUILD) -o ../builds/$(BINARY_NAME)-$(VERSION)-darwin-amd64 -v
	cd cmd; env GOOS=linux GOARCH=amd64 $(GOBUILD) -o ../builds/$(BINARY_NAME)-$(VERSION)-linux-amd64 -v
	cd cmd; env GOOS=windows GOARCH=amd64 $(GOBUILD) -o ../builds/$(BINARY_NAME)-$(VERSION)-windows-amd64.exe -v

install: build ## build and install the binary
	cp builds/$(BINARY_NAME) $(GOPATH)/bin/$(BINARY_NAME)
	# md5 builds/$(BINARY_NAME) || md5sum builds/$(BINARY_NAME)

deps: ## Install dependencies
	$(GOGET) ./...
	# $(GOGET) github.com/inconshreveable/mousetrap
	# $(GOGET) github.com/konsorten/go-windows-terminal-sequences

# docker: image ## Build docker image and upload to docker hub
# 	docker login

# image: clean ## Build docker image
# 	docker build -t $(BINARY_NAME) .

test: deps ## Run unit test
	$(GOTEST) -v ./...

clean: ## Remove files created by the build
	$(GOCLEAN)
	rm -fr builds

# run: ## Run server docker image
# 	docker run -it --rm -p 8080:8080 urlinsane

doc: ## Go documentation
	$(GODOC) -http=:6060

# login-gcr: ## docker login to GCR
# 	docker login -u oauth2accesstoken -p "$(shell gcloud auth print-access-token)" https://$(GCR_HOST)

# push-gcr: login-gcr image ## Push build to Google Container Registry
# 	docker tag $(BINARY_NAME) $(GCR_HOST)/$(GCP_PROJECT_ID)/$(BINARY_NAME)
# 	docker push $(GCR_HOST)/$(GCP_PROJECT_ID)/$(BINARY_NAME)

# deploy-gc-app-engine: push-gcr ## Deploy api service to Google Cloud AppEngine
# 	gcloud app deploy --quiet --image-url $(GCR_HOST)/$(GCP_PROJECT_ID)/$(BINARY_NAME)

