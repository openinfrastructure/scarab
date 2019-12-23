# The target device to install to
TARGET_HOST ?= gw.home
VERSION ?= $(shell git describe --tags 2> /dev/null | cut -d- -f1 || echo 0.0.0)
BUILD ?= $(shell git describe --tags 2> /dev/null | cut -s -d- -f2-)
PROJECTNAME ?= $(shell basename "$(PWD)")
BUILD_TIME ?= $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
MODULE := github.com/openinfrastructure/$(PROJECTNAME)

# Go related variables.
LDFLAGS=-ldflags "-X=$(MODULE)/common/scarab.BuildVersion=$(VERSION) -X=$(MODULE)/common/scarab.Build=$(BUILD) -X=$(MODULE)/common/scarab.BuildTime=$(BUILD_TIME)"

build:
	go build $(LDFLAGS) -o bin/$(PROJECTNAME) main.go

build.mips:
	GOOS=linux GOARCH=mips64 go build $(LDFLAGS) -o bin.mips64/$(PROJECTNAME) main.go

lint:
	golangci-lint run ./...

test:
	go test -coverprofile cover.out ./...

cover: test
	go tool cover -html=cover.out -o cover.html

check: lint test cover

fmt:
	go fmt ./...

install:
	envsubst < scarab.yaml.sample > scarab.yaml
	ssh -q ubnt@$(TARGET_HOST) 'mkdir -p /config/scarab/bin /config/scarab/bin.mips64 /config/scarab /config/scripts/ppp/ip-up.d'
	scp -q scarab.yaml ubnt@$(TARGET_HOST):/config/scarab/scarab.yaml
	scp -q bin.mips64/$(PROJECTNAME) ubnt@$(TARGET_HOST):/config/scarab/bin.mips64/$(PROJECTNAME)
	scp -q scripts/$(PROJECTNAME) ubnt@$(TARGET_HOST):/config/scarab/bin/$(PROJECTNAME)
	scp -q scripts/ip-up.d/scarab ubnt@$(TARGET_HOST):/config/scripts/ppp/ip-up.d/scarab
