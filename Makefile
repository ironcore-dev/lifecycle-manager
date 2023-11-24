CONTROLLER_TOOLS_VERSION ?= v0.13.0
ADDLICENSE_VERSION ?= v1.1.1
GOLANGCI_LINT_VERSION ?= v1.55.2
COPYRIGHT ?= "T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved"

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: lint
lint: golangci-lint
	$(GOLANGCI_LINT) run ./...

.PHONY: generate
generate: controller-gen ## Generate code containing DeepCopy, DeepCopyInto, and DeepCopyObject method implementations.
	$(CONTROLLER_GEN) object:headerFile="hack/copyright.txt" paths="./..."

.PHONY: add-license
add-license: addlicense ## Add license header to all .go files in project
	@find . -name '*.go' -exec $(ADDLICENSE) -c $(COPYRIGHT) {} +

.PHONY: check-license
check-license: addlicense ## Check license header presence in all .go files in project
	@find . -name '*.go' -exec $(ADDLICENSE) -check {} +

### AUXILIARY ###
ADDLICENSE = $(shell pwd)/bin/addlicense
.PHONY: addlicense
addlicense: ##Download addlicense to local project's bin folder
	$(call go-get-tool,$(ADDLICENSE),github.com/google/addlicense@$(ADDLICENSE_VERSION))

CONTROLLER_GEN = $(shell pwd)/bin/controller-gen
.PHONY: controller-gen
controller-gen: ## Download controller-gen locally if necessary.
	$(call go-get-tool,$(CONTROLLER_GEN),sigs.k8s.io/controller-tools/cmd/controller-gen@$(CONTROLLER_TOOLS_VERSION))

GOLANGCI_LINT = $(shell pwd)/bin/golangci-lint
.PHONY: golangci-lint
golangci-lint:
	$(call go-get-tool,$(GOLANGCI_LINT),github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION))

# go-get-tool will 'go install' any package $2 and install it to $1.
PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
define go-get-tool
@[ -f $(1) ] || { \
set -e ; \
GOBIN=$(PROJECT_DIR)/bin go install $(2) ; \
}
endef
