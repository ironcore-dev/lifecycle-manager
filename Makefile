.PHONY: fmt
fmt: goimports
	go fmt ./...
	$(GOIMPORTS) -w .

.PHONY: vet
vet:
	go vet ./...

.PHONY: lint
lint: golangci-lint
	$(GOLANGCI_LINT) run

.PHONY: lint-fix
lint-fix: golangci-lint
	$(GOLANGCI_LINT) run --fix

.PHONY: manifests
manifests: controller-gen ## Generate WebhookConfiguration, ClusterRole and CustomResourceDefinition objects.
	$(CONTROLLER_GEN) rbac:roleName=manager-role crd webhook paths="./..." output:crd:artifacts:config=config/crd/bases

.PHONY: generate
generate: code-gen

.PHONY: add-license
add-license: addlicense ## Add license header to all .go files in project
	@find . -name '*.go' -exec $(ADDLICENSE) -f hack/license-header.txt {} +

.PHONY: check-license
check-license: addlicense ## Check license header presence in all .go files in project
	@find . -name '*.go' -exec $(ADDLICENSE) -check -c 'IronCore authors' {} +

.PHONY: test-controllers
test-controllers:
	go test ./internal/... -coverprofile cover.out

.PHONY: test-integration
test-integration: envtest
	KUBEBUILDER_ASSETS="$(shell $(ENVTEST) use $(ENVTEST_K8S_VERSION) -p path)" go test ./integrationtests/... -coverprofile cover.out

.PHONY: test
test: test-controllers test-integration

.PHONY: check
check: vet lint check-license test

.PHONY: format
format: generate manifests add-license fmt lint-fix

.PHONY: docs
docs: gen-crd-api-reference-docs ## Run go generate to generate API reference documentation.
	$(GEN_CRD_API_REFERENCE_DOCS) -api-dir ./api/lifecycle/v1alpha1 -config ./hack/api-reference/config.json -template-dir ./hack/api-reference/template -out-file ./docs/api-reference/lifecycle.md

### AUXILIARY ###
LOCAL_BIN ?= $(shell pwd)/bin
$(LOCAL_BIN):
	mkdir -p $(LOCAL_BIN)

## Tools locations
ADDLICENSE ?= $(LOCAL_BIN)/addlicense
CONTROLLER_GEN ?= $(LOCAL_BIN)/controller-gen
GOLANGCI_LINT ?= $(LOCAL_BIN)/golangci-lint
GOIMPORTS ?= $(LOCAL_BIN)/goimports
PROTOC_GEN_GOGO ?= $(LOCAL_BIN)/protoc-gen-gogo
ENVTEST ?= $(LOCAL_BIN)/setup-envtest
DEEPCOPY_GEN ?= $(LOCAL_BIN)/deepcopy-gen
CLIENT_GEN ?= $(LOCAL_BIN)/client-gen
LISTER_GEN ?= $(LOCAL_BIN)/lister-gen
INFORMER_GEN ?= $(LOCAL_BIN)/informer-gen
DEFAULTER_GEN ?= $(LOCAL_BIN)/defaulter-gen
CONVERSION_GEN ?= $(LOCAL_BIN)/conversion-gen
OPENAPI_GEN ?= $(LOCAL_BIN)/openapi-gen
APPLYCONFIGURATION_GEN ?= $(LOCAL_BIN)/applyconfiguration-gen
GO_TO_PROTOBUF ?= $(LOCAL_BIN)/go-to-protobuf
MODELS_SCHEMA ?= $(LOCAL_BIN)/models-schema
VGOPATH ?= $(LOCAL_BIN)/vgopath
GEN_CRD_API_REFERENCE_DOCS ?= $(LOCAL_BIN)/gen-crd-api-reference-docs

## Tools versions
ADDLICENSE_VERSION ?= v1.1.1
CONTROLLER_GEN_VERSION ?= v0.13.0
GOLANGCI_LINT_VERSION ?= v1.55.2
GOIMPORTS_VERSION ?= v0.16.1
PROTOC_GEN_GOGO_VERSION ?= v1.3.2
ENVTEST_K8S_VERSION ?= 1.28.3
CODE_GENERATOR_VERSION ?= v0.28.3
VGOPATH_VERSION ?= v0.1.3
GEN_CRD_API_REFERENCE_DOCS_VERSION ?= v0.3.0

.PHONY: code-gen
code-gen: vgopath goimports go-to-protobuf deepcopy-gen protoc-gen-gogo
	@VGOPATH=$(VGOPATH) \
	DEEPCOPY_GEN=$(DEEPCOPY_GEN) \
	GO_TO_PROTOBUF=$(GO_TO_PROTOBUF) \
	PROTOC_GEN_GOGO=$(PROTOC_GEN_GOGO) \
	./hack/generate.sh

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

.PHONY: goimports
goimports: $(GOIMPORTS)
$(GOIMPORTS): $(LOCAL_BIN)
	@test -s $(GOIMPORTS) || GOBIN=$(LOCAL_BIN) go install golang.org/x/tools/cmd/goimports@$(GOIMPORTS_VERSION)

.PHONY: protoc-gen-gogo
protoc-gen-gogo: $(PROTOC_GEN_GOGO)
$(PROTOC_GEN_GOGO): $(LOCAL_BIN)
	@test -s $(PROTOC_GEN_GOGO) || GOBIN=$(LOCAL_BIN) go install github.com/gogo/protobuf/protoc-gen-gogo@$(PROTOC_GEN_GOGO_VERSION)

.PHONY: envtest
envtest: $(ENVTEST) ## Download envtest-setup locally if necessary.
$(ENVTEST): $(LOCAL_BIN)
	@test -s $(ENVTEST) || GOBIN=$(LOCAL_BIN) go install sigs.k8s.io/controller-runtime/tools/setup-envtest@latest

.PHONY: vgopath
vgopath: $(VGOPATH)
$(VGOPATH): $(LOCAL_BIN)
	@test -s $(VGOPATH) || GOBIN=$(LOCAL_BIN) go install github.com/ironcore-dev/vgopath@$(VGOPATH_VERSION)

.PHONY: deepcopy-gen
deepcopy-gen: $(DEEPCOPY_GEN)
$(DEEPCOPY_GEN): $(LOCAL_BIN)
	@test -s $(DEEPCOPY_GEN) || GOBIN=$(LOCAL_BIN) go install k8s.io/code-generator/cmd/deepcopy-gen@$(CODE_GENERATOR_VERSION)

.PHONY: go-to-protobuf
go-to-protobuf: $(GO_TO_PROTOBUF)
$(GO_TO_PROTOBUF): $(LOCAL_BIN)
	@test -s $(GO_TO_PROTOBUF) || GOBIN=$(LOCAL_BIN) go install k8s.io/code-generator/cmd/go-to-protobuf@$(CODE_GENERATOR_VERSION)

.PHONY: gen-crd-api-reference-docs
gen-crd-api-reference-docs: $(GEN_CRD_API_REFERENCE_DOCS) ## Download gen-crd-api-reference-docs locally if necessary.
$(GEN_CRD_API_REFERENCE_DOCS): $(LOCAL_BIN)
	test -s $(GEN_CRD_API_REFERENCE_DOCS) || GOBIN=$(LOCAL_BIN) go install github.com/ahmetb/gen-crd-api-reference-docs@$(GEN_CRD_API_REFERENCE_DOCS_VERSION)

