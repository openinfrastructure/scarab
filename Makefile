VERSION := $(shell git describe --tags 2> /dev/null | cut -d- -f1 || echo 0.0.0)
BUILD := $(shell git describe --tags 2> /dev/null | cut -s -d- -f2-)
PROJECTNAME := $(shell basename "$(PWD)")
BUILD_TIME := $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
MODULE := github.com/jeffmccune/$(PROJECTNAME)

# Go related variables.
GOBASE := $(shell pwd)
GOBIN := $(GOBASE)/bin
GOFILES := $(wildcard *.go)
LDFLAGS=-ldflags "-X=$(MODULE)/common/scarab.BuildVersion=$(VERSION) -X=$(MODULE)/common/scarab.Build=$(BUILD) -X=$(MODULE)/common/scarab.BuildTime=$(BUILD_TIME)"

build:
	go build $(LDFLAGS) -o bin/$(PROJECTNAME) main.go