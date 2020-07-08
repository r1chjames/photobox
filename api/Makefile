-include .env

BUILD := $(shell git rev-parse --short HEAD)
PROJECTNAME := $(shell basename "$(PWD)")

# Go related variables.
GOBASE := $(shell pwd)
GOPATH := $(GOBASE)/vendor:$(GOBASE)
GOBIN := $(GOBASE)/bin
GOFILES := $(wildcard *.go)

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent
# Redirect error output to a file, so we can show it in development mode.
STDERR := /tmp/.$(PROJECTNAME)-stderr.txt

.PHONY: build clean deploy

## install: Install missing dependencies. Runs `go get` internally. e.g; make install get=github.com/foo/bar
install: go-get

deploy:

build: clean install

## compile: Compile the binary.
compile:
	@-touch $(STDERR)
	@-rm $(STDERR)
	@-$(MAKE) -s go-build 2> $(STDERR)
	@cat $(STDERR) | sed -e '1s/.*/\nError:\n/'  | sed 's/make\[.*/ /' | sed "/^/s/^/     /" 1>&2

## exec: Run given command, wrapped with custom GOPATH. e.g; make exec run="go test ./..."
exec:
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) $(run)

## clean: Clean build files. Runs `go clean` internally.
clean:
	@-rm $(GOBIN)/$(PROJECTNAME) 2> /dev/null
	@-$(MAKE) go-clean

go-build:
	@echo "  >  Building binary..."
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) CGO_ENABLED=0 go build $(LDFLAGS) -o $(GOBIN)/$(PROJECTNAME) $(GOFILES)

go-generate:
	@echo "  >  Generating dependency files..."
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go generate $(generate)

go-get:
	@echo "  >  Checking if there are any missing dependencies..."
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go mod download

go-install:
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go install $(GOFILES)

go-clean:
	@echo "  >  Cleaning build cache"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go clean

go-test:
	@echo "Running tests"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go test ./...

docker-build:
	docker build -t registry.gitlab.com/r1chjames/photobox .

docker-push: docker-build
	docker push registry.gitlab.com/r1chjames/photobox

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo