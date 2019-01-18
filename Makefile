# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=urlinsane

VERSION=0.4.0

all: build hash
hash:
	md5 builds/$(BINARY_NAME) | md5sum builds/$(BINARY_NAME)
build:
	mkdir -p builds
	cd cmd; $(GOBUILD) -o ../builds/$(BINARY_NAME) -v
test:
	$(GOTEST) -v ./...
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
run:
	cd cmd; $(GOBUILD) -o ../$(BINARY_NAME) -v
	./$(BINARY_NAME)
deps:
	$(GOGET) github.com/rangertaha/urlinsane
versions:
	mkdir -p builds
	cd cmd; env GOOS=darwin GOARCH=amd64 $(GOBUILD) -o ../builds/$(BINARY_NAME)-$(VERSION)-darwin-amd64 -v
	cd cmd; env GOOS=linux GOARCH=amd64 $(GOBUILD) -o ../builds/$(BINARY_NAME)-$(VERSION)-linux-amd64 -v
	cd cmd; env GOOS=windows GOARCH=amd64 $(GOBUILD) -o ../builds/$(BINARY_NAME)-$(VERSION)-windows-amd64.exe -v
