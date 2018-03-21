## Project Vars ##########################################################
GO_VARS = ENABLE_CGO=0 GOOS=darwin GOARCH=amd64
GO ?= go
GIT ?= git
COMMIT := $(shell $(GIT) rev-parse HEAD)
VERSION ?= $(shell $(GIT) describe --tags ${COMMIT} 2> /dev/null || echo "$(COMMIT)")
BUILD_TIME := $(shell LANG=en_US date +"%F_%T_%z")
ROOT := gitlab.com/kanalbot/ershad
ROOT_DIRECTORY := $(shell pwd)
LD_FLAGS := -X $(ROOT).Version=$(VERSION) -X $(ROOT).Commit=$(COMMIT) -X $(ROOT).BuildTime=$(BUILD_TIME) -X $(ROOT).Title=ershadd

.PHONY: help clean

ershadd: *.go */*.go */*/*.go Gopkg.lock
	$(GO_VARS) $(GO) build -i -o="ershadd" -ldflags="$(LD_FLAGS)" $(ROOT)/cmd/ershad

clean:
	rm -rf ershadd

help:
	@echo "Please use \`make <ROOT>' where <ROOT> is one of"
	@echo "  ershadd     to build the main binary for current platform"
	@echo "  clean             to remove generated files"
