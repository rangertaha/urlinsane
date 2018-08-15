# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=urlinsane

all: test build hash
hash:
	md5 $(BINARY_NAME)
build:
	cd cmd; $(GOBUILD) -o ../$(BINARY_NAME) -v
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
