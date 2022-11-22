GO=go
GOTEST=$(GO) test
GOVET=$(GO) vet
OUT_DIR := $(if $(OUT_DIR),$(OUT_DIR),./bin)
MOCKS_DIR=./internal/mocks

export GO111MODULE := on

GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
CYAN   := $(shell tput -Txterm setaf 6)
RESET  := $(shell tput -Txterm sgr0)

M = $(shell printf "\033[32;1m▶▶▶\033[0m")
M2 = $(shell printf "\033[32;1m▶▶▶▶▶▶\033[0m")

.PHONY: all tools mocks test build vendor

all: help

## Build:
tools: ## Install all the tools needed to build artifacts
	@echo '$(M) downloading tools…'
	@echo "$(M2) Installing mockery…"
	go install github.com/vektra/mockery/v2@v2.15.0

mocks: ## Generate mocks for Golang interfaces
	@echo '$(M) generating mocks for Golang interfaces…'
	rm -rf $(MOCKS_DIR)
	mockery --all --keeptree --output=$(MOCKS_DIR)

build: ## Build your project and put the output binary in out/bin/
	@echo '$(M) building the artifacts…'
	mkdir -p $(OUT_DIR)
	$(GO) build -o $(OUT_DIR) ./cmd/...

clean: ## Remove build related file
	@echo '$(M) removing the artifacts…'
	rm -rf $(OUT_DIR)

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