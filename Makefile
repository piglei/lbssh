# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
VERSION=0.0.2

LDFLAGS=-ldflags "-X github.com/piglei/lbssh/pkg/version.version=$(VERSION) \
-X github.com/piglei/lbssh/pkg/version.gitCommit=`git rev-parse HEAD` \
-X github.com/piglei/lbssh/pkg/version.buildDate=`date -u +'%Y-%m-%dT%H:%M:%SZ'`"

.PHONY: build test clean env-test

all: test build
build:
	$(GOBUILD) $(LDFLAGS) -v -o build/lbssh
test:
	$(GOTEST) -v ./... -cover
clean:
	$(GOCLEAN)
build-all:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -v -o build/lbssh_darwin_amd64
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -v -o build/lbssh_linux_amd64

