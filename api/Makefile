-include .env

BUILD := $(shell git rev-parse --short HEAD)
PROJECTNAME := $(shell basename "$(PWD)")

# Go related variables.
GOBASE := $(shell pwd)
# GOPATH := $(GOBASE)/vendor:$(GOBASE)
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

build: clean install compile

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
	@GOPATH=$(GOPATH) GO111MODULE=on GOBIN=$(GOBIN) CGO_ENABLED=0 go build $(LDFLAGS) -o $(GOBIN)/$(PROJECTNAME) $(GOFILES)

go-generate:
	@echo "  >  Generating dependency files..."
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go generate $(generate)

go-get:
	@echo "  >  Checking if there are any missing dependencies..."
	@GOBIN=$(GOBIN) go mod vendor

go-install:
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go install $(GOFILES)

go-clean:
	@echo "  >  Cleaning build cache"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go clean

## test: Run unit tests only (fast)
go-test:
	@echo "  >  Running unit tests..."
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go test ./internal/... -short

run-db:
	@echo "Starting local DB"
	./scripts/db.sh

## test-integration-up: Start test database for integration tests
test-integration-up:
	@echo "  >  Starting test database..."
	@docker-compose -f test/docker-compose.test.yml up -d
	@echo "  >  Waiting for database to be ready..."
	@sleep 3

## test-integration-down: Stop test database
test-integration-down:
	@echo "  >  Stopping test database..."
	@docker-compose -f test/docker-compose.test.yml down

## test-integration-clean: Stop and remove test database volumes
test-integration-clean:
	@echo "  >  Cleaning test database..."
	@docker-compose -f test/docker-compose.test.yml down -v

## test-integration: Run integration tests (requires test DB running)
test-integration:
	@echo "  >  Running integration tests..."
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go test -tags=integration ./test/integration/... -v -timeout 30m

## test-integration-full: Start DB, run integration tests, stop DB
test-integration-full: test-integration-up
	@echo "  >  Running integration tests..."
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go test -tags=integration ./test/integration/... -v -timeout 30m || ($(MAKE) test-integration-down && exit 1)
	@$(MAKE) test-integration-down

## test-all: Run both unit and integration tests
test-all: go-test test-integration-full

## test-coverage: Generate test coverage report
test-coverage:
	@echo "  >  Generating coverage report..."
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go test ./internal/... -coverprofile=coverage.out -covermode=atomic
	@go tool cover -html=coverage.out -o coverage.html
	@echo "  >  Coverage report: coverage.html"

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo