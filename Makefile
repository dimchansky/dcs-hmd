GO=go
GOTEST=$(GO) test
GOVET=$(GO) vet
PACKAGE_NAME := $(shell $(GO) list -m)
OUT_DIR := $(if $(OUT_DIR),$(OUT_DIR),./bin)
MOCKS_DIR=./internal/mocks

BENCH ?=.
BENCH_FLAGS ?= -benchmem
export GO111MODULE := on

VERSION ?= vlatest
COMMIT ?= $(shell git rev-parse HEAD)
BUILD_TIME ?= $(shell TZ=UTC date -u '+%Y-%m-%dT%H:%M:%SZ')

GO_LDFLAGS := '-X "$(PACKAGE_NAME)/cmd.Version=$(VERSION)" -X "$(PACKAGE_NAME)/cmd.BuildTime=$(BUILD_TIME)" -X "$(PACKAGE_NAME)/cmd.GitHash=$(COMMIT)"'

GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
CYAN   := $(shell tput -Txterm setaf 6)
RESET  := $(shell tput -Txterm sgr0)

M = $(shell printf "\033[32;1m▶▶▶\033[0m")
M2 = $(shell printf "\033[32;1m▶▶▶▶▶▶\033[0m")

.PHONY: all tools mocks test bench build vendor

all: help

## Build:
tools: ## Install all the tools needed to build artifacts
	@echo '$(M) downloading tools…'
	@echo "$(M2) Installing mockery…"
	go install github.com/vektra/mockery/v2@v2.22.1
	@echo "$(M2) Installing goimports..."
	go install golang.org/x/tools/cmd/goimports@v0.3.0
	@echo "$(M2) Installing golangci-lint..."
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.50.1

mocks: ## Generate mocks for Golang interfaces
	@echo '$(M) generating mocks for Golang interfaces…'
	rm -rf $(MOCKS_DIR)
	mockery --all --keeptree --output=$(MOCKS_DIR)

build: ## Build your project and put the output binary in out/bin/
	@echo '$(M) building the artifacts…'
	mkdir -p $(OUT_DIR)
	$(GO) build --ldflags=$(GO_LDFLAGS) -o $(OUT_DIR) ./cmd/...

clean: ## Remove build related file
	@echo '$(M) removing the artifacts…'
	rm -rf $(OUT_DIR)

## Test:
test: ## Run the tests of the project
	@echo '$(M) running tests…'
	$(GOTEST) -v -race ./...

bench: ## Run the benchmarks of the project
	@echo '$(M) running benchmarks…'
	$(GOTEST) -bench=$(BENCH) -run="^$$" $(BENCH_FLAGS) ./...

## Lint:
lint: ## Use golintci-lint on your project
	@echo '$(M) running golangci-lint…'
	golangci-lint run

## Help:
help: ## Show this help.
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "    ${YELLOW}%-20s${GREEN}%s${RESET}\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  ${CYAN}%s${RESET}\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)