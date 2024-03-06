# Build the manager binary
FROM golang:1.22.0 as builder

ARG GOARCH

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum

COPY hack/ hack/
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN mkdir -p -m 0600 ~/.ssh \
  && ssh-keyscan -t rsa github.com >> ~/.ssh/known_hosts \
  && go mod download


# Copy the go source
COPY cmd/lifecycle-controller-manager/main.go cmd/lifecycle-controller-manager/main.go
COPY api/ api/
COPY lcmi/ lcmi/
COPY internal/ internal/
COPY util/ util/

# Build
# the GOARCH has not a default value to allow the binary be built according to the host where the command
# was called. For example, if we call make docker-build in a local env which has the Apple Silicon M1 SO
# the docker BUILDPLATFORM arg will be linux/arm64 when for Apple x86 it will be linux/amd64. Therefore,
# by leaving it empty we can ensure that the container and binary shipped on it will have the same platform.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${GOARCH} go build -a -o manager cmd/lifecycle-controller-manager/main.go

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/manager .
USER 65532:65532

ENTRYPOINT ["/manager"]
