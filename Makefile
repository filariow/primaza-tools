##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

.PHONY: help
help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-30s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

## Common

PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
LOCALBIN := bin

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# Setting SHELL to bash allows bash commands to be executed by recipes.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

OUTPUT_DIR ?= $(PROJECT_DIR)/out
$(OUTPUT_DIR):
	mkdir -p $(OUTPUT_DIR)
HACK_DIR ?= $(PROJECT_DIR)/hack
GOCACHE ?= $(OUTPUT_DIR)/.gocache
GOFLAGS ?=
GO ?= GOCACHE=$(GOCACHE) GOFLAGS="$(GOFLAGS)" go

## Build Dependencies

## Location to install dependencies to
LOCALBIN ?= $(PROJECT_DIR)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

##@ Development
.PHONY: test
test:
	$(GO) test ./...

##@ Build
build: ## build the cli
	$(GO) build -ldflags="-s -w" -trimpath -o out/primaza-adm ./cmd/primaza-adm/main.go

##@ Linters

GOLANGCI_LINT=$(LOCALBIN)/golangci-lint
GOLANGCI_LINT_VERSION ?= v1.51.2

YAMLLINT_VERSION ?= 1.28.0

SHELLCHECK=$(LOCALBIN)/shellcheck
SHELLCHECK_VERSION ?= v0.9.0

GO_LINT_CMD = GOFLAGS="$(GOFLAGS)" GOGC=30 GOCACHE=$(GOCACHE) $(GOLANGCI_LINT) run

.PHONY: fmt
fmt: ## Run go fmt against code.
	$(GO) fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	$(GO) vet ./...

.PHONY: lint-go
lint-go: $(GOLANGCI_LINT) fmt vet ## Checks Go code
	$(GO_LINT_CMD)

$(GOLANGCI_LINT):
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(LOCALBIN) $(GOLANGCI_LINT_VERSION)

.PHONY: shellcheck
shellcheck: $(SHELLCHECK) ## Download shellcheck locally if necessary.
$(SHELLCHECK): $(OUTPUT_DIR)
ifeq (,$(wildcard $(SHELLCHECK)))
ifeq (,$(shell which shellcheck 2>/dev/null))
	@{ \
	set -e ;\
	mkdir -p $(dir $(SHELLCHECK)) ;\
	OS=$(shell go env GOOS) && ARCH=$(shell go env GOARCH | sed -e 's,amd64,x86_64,g') && \
	curl -Lo $(OUTPUT_DIR)/shellcheck.tar.xz https://github.com/koalaman/shellcheck/releases/download/$(SHELLCHECK_VERSION)/shellcheck-$(SHELLCHECK_VERSION).$${OS}.$${ARCH}.tar.xz ;\
	tar --directory $(OUTPUT_DIR) -xvf $(OUTPUT_DIR)/shellcheck.tar.xz ;\
	find $(OUTPUT_DIR) -name shellcheck -exec cp {} $(SHELLCHECK) \; ;\
	chmod +x $(SHELLCHECK) ;\
	}
else
SHELLCHECK = $(shell which shellcheck)
endif
endif

.PHONY: clean
clean: ## Removes temp directories
	-rm -rf ${V_FLAG} $(OUTPUT_DIR)
