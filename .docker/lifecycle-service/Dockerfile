# Build the manager binary
FROM golang:1.22.3 as builder

ARG GOARCH

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

# Copy the go source
COPY cmd/lifecycle-service/ cmd/lifecycle-service
COPY api/ api/
COPY clientgo/applyconfiguration clientgo/applyconfiguration
COPY clientgo/lifecycle clientgo/lifecycle
COPY clientgo/connectrpc clientgo/connectrpc
COPY internal/service internal/service
COPY internal/util internal/util

# Build
# the GOARCH has not a default value to allow the binary be built according to the host where the command
# was called. For example, if we call make docker-build in a local env which has the Apple Silicon M1 SO
# the docker BUILDPLATFORM arg will be linux/arm64 when for Apple x86 it will be linux/amd64. Therefore,
# by leaving it empty we can ensure that the container and binary shipped on it will have the same platform.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${GOARCH} go build -a -o manager cmd/lifecycle-service/main.go

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/manager .
USER 65532:65532

ENTRYPOINT ["/manager"]
