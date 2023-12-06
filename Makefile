

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
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."

.PHONY: add-license
add-license: addlicense ## Add license header to all .go files in project
	@find . -name '*.go' -exec $(ADDLICENSE) -f hack/license-header.txt {} +

.PHONY: check-license
check-license: addlicense ## Check license header presence in all .go files in project
	@find . -name '*.go' -exec $(ADDLICENSE) -check -c 'IronCore authors' {} +

.PHONY: test
test:
	go test ./... -coverprofile cover.out

### AUXILIARY ###
LOCAL_BIN ?= $(shell pwd)/bin
$(LOCAL_BIN):
	mkdir -p $(LOCAL_BIN)

## Tools locations
ADDLICENSE = $(LOCAL_BIN)/addlicense
CONTROLLER_GEN = $(LOCAL_BIN)/controller-gen
GOLANGCI_LINT = $(LOCAL_BIN)/golangci-lint

## Tools versions
CONTROLLER_GEN_VERSION ?= v0.13.0
ADDLICENSE_VERSION ?= v1.1.1
GOLANGCI_LINT_VERSION ?= v1.55.2

.PHONY: addlicense
addlicense: ##Download addlicense to local project's bin folder
	@test -s $(ADDLICENSE) || GOBIN=$(LOCAL_BIN) go install github.com/google/addlicense@$(ADDLICENSE_VERSION)

.PHONY: controller-gen
controller-gen: ## Download controller-gen locally if necessary.
	@test -s $(CONTROLLER_GEN) && $(CONTROLLER_GEN) --version | grep -q $(CONTROLLER_GEN_VERSION) || \
	GOBIN=$(LOCAL_BIN) go install sigs.k8s.io/controller-tools/cmd/controller-gen@$(CONTROLLER_GEN_VERSION)

.PHONY: golangci-lint
golangci-lint:
	@test -s $(GOLANGCI_LINT) && $(GOLANGCI_LINT) --version | grep -q $(GOLANGCI_LINT_VERSION) || \
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)
