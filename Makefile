

.PHONY: fmt
fmt: buf
	go fmt ./...
	$(BUF) format -w

.PHONY: vet
vet:
	go vet ./...

.PHONY: lint
lint: golangci-lint buf
	$(GOLANGCI_LINT) run ./...
	$(BUF) lint

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

.PHONY: proto
proto: buf protoc-gen-gogotypes
	$(BUF) generate


### AUXILIARY ###
LOCAL_BIN ?= $(shell pwd)/bin
$(LOCAL_BIN):
	mkdir -p $(LOCAL_BIN)

## Tools locations
ADDLICENSE ?= $(LOCAL_BIN)/addlicense
CONTROLLER_GEN ?= $(LOCAL_BIN)/controller-gen
GOLANGCI_LINT ?= $(LOCAL_BIN)/golangci-lint
PROTOC_GEN_GOGO_TYPES ?= $(LOCAL_BIN)/protoc-gen-gogotypes
BUF ?= $(LOCAL_BIN)/buf

## Tools versions
CONTROLLER_GEN_VERSION ?= v0.13.0
ADDLICENSE_VERSION ?= v1.1.1
GOLANGCI_LINT_VERSION ?= v1.55.2
PROTOC_GEN_GOGO_TYPES_VERSION ?= v1.3.2
BUF_VERSION ?= v1.28.1

.PHONY: addlicense
addlicense: $(ADDLICENSE)
$(ADDLICENSE): $(LOCAL_BIN)
	@test -s $(ADDLICENSE) || GOBIN=$(LOCAL_BIN) go install github.com/google/addlicense@$(ADDLICENSE_VERSION)

.PHONY: controller-gen
controller-gen: $(CONTROLLER_GEN)
$(CONTROLLER_GEN): $(LOCAL_BIN)
	@test -s $(CONTROLLER_GEN) && $(CONTROLLER_GEN) --version | grep -q $(CONTROLLER_GEN_VERSION) || \
	GOBIN=$(LOCAL_BIN) go install sigs.k8s.io/controller-tools/cmd/controller-gen@$(CONTROLLER_GEN_VERSION)

.PHONY: golangci-lint
golangci-lint: $(GOLANGCI_LINT)
$(GOLANGCI_LINT): $(LOCAL_BIN)
	@test -s $(GOLANGCI_LINT) && $(GOLANGCI_LINT) --version | grep -q $(GOLANGCI_LINT_VERSION) || \
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

.PHONY: protoc-gen-gogotypes
protoc-gen-gogotypes: $(PROTOC_GEN_GOGO_TYPES)
$(PROTOC_GEN_GOGO_TYPES): $(LOCAL_BIN)
	@test -s $(PROTOC_GEN_GOGO_TYPES) || GOBIN=$(LOCAL_BIN) go install github.com/gogo/protobuf/protoc-gen-gogotypes@$(PROTOC_GEN_GOGO_TYPES_VERSION)

.PHONY: buf
buf: $(BUF)
$(BUF): $(LOCAL_BIN)
	@test -s $(BUF) || GOBIN=$(LOCAL_BIN) go install github.com/bufbuild/buf/cmd/buf@$(BUF_VERSION)