PKG_LIST := $(shell go list ./... | grep -v /vendor/)

.SILENT: ;
.PHONY: all

all: build

lint: ## Lint the files
	golint -set_exit_status $(PKG_LIST)

test: ## Run unit tests
	go fmt $(PKG_LIST)
	go vet $(PKG_LIST)
	go test -race -timeout 30s -cover -v -count 1 $(PKG_LIST)

bench: ## Run benchmark
	go test -bench . $(PKG_LIST)

msan: ## Run memory sanitizer
	go test -msan $(PKG_LIST)

build: ## Build the binary file
	go build -v $(PKG_LIST)

cover: ## Code coverage
	go test -coverprofile=cover.out $(PKG_LIST)

clean: ## Remove previous build
	rm -f cover.out
	go clean

help: ## Display this help screen
	grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
