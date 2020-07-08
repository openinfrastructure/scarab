VERSION := $(shell git describe --tags 2> /dev/null | cut -d- -f1 || echo 0.0.0)
BUILD := $(shell git describe --tags 2> /dev/null | cut -s -d- -f2-)
PROJECTNAME := $(shell basename "$(PWD)")
BUILD_TIME := $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
MODULE := github.com/openinfrastructure/$(PROJECTNAME)

# Go related variables.
LDFLAGS=-ldflags "-X=$(MODULE)/common/scarab.BuildVersion=$(VERSION) -X=$(MODULE)/common/scarab.Build=$(BUILD) -X=$(MODULE)/common/scarab.BuildTime=$(BUILD_TIME)"

build:
	go build $(LDFLAGS) -o bin/$(PROJECTNAME) main.go

# Build for the EdgeRouter-4
build.er4:
	GOOS=linux GOARCH=mips64 go build $(LDFLAGS) -o bin/$(PROJECTNAME).mips64 main.go

install.er4: build.er4
	scp bin/$(PROJECTNAME).mips64 ubnt@gw:/config/scarab/bin.mips64/scarab

lint:
	golangci-lint run ./...

test:
	go test -coverprofile cover.out ./...

cover: test
	go tool cover -html=cover.out -o cover.html

check: lint test cover

fmt:
	go fmt ./...

# dlv debug -- dns update --dnszone=jeff --rrnames="*.jeff.ois.run."
# Type 'help' for list of commands.
# (dlv) break cmd.debug
# Breakpoint 1 set at 0x1a17cd3 for github.com/openinfrastructure/scarab/cmd.debug() ./cmd/dns_update.go:129
# (dlv) continue
# Version string empty
# > github.com/openinfrastructure/scarab/cmd.debug() ./cmd/dns_update.go:129 (hits goroutine(1):1 total:1) (PC: 0x1a17cd3)
#    124:         }
#    125:
#    126:         log.Println(c)
#    127: }
#    128:
# => 129: func debug(c *dns.Change) bool {
#    130:         log.Println(c)
#    131:         return false
#    132: }
# (dlv)
debug:
	dlv debug -- dns update --dnszone=jeff --rrnames="*.jeff.ois.run."
