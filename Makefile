# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

.PHONY: build test minikube clean linux env-test

all: test build
build:
	$(GOBUILD) -v -o build/lbssh
test:
	$(GOTEST) -v ./... -cover
clean:
	$(GOCLEAN)
build-all:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -v -o build/lbssh_darwin_amd64
	GOOS=linux GOARCH=amd64 $(GOBUILD) -v -o build/lbssh_linux_amd64

