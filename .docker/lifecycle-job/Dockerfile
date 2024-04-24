# Build the manager binary
FROM golang:1.22.2 as builder

ARG GOARCH

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

# Copy the go source
COPY cmd/lifecycle-job cmd/lifecycle-job
COPY api/ api/
COPY clientgo/lifecycle clientgo/lifecycle
COPY clientgo/connectrpc clientgo/connectrpc

# Build
# the GOARCH has not a default value to allow the binary be built according to the host where the command
# was called. For example, if we call make docker-build in a local env which has the Apple Silicon M1 SO
# the docker BUILDPLATFORM arg will be linux/arm64 when for Apple x86 it will be linux/amd64. Therefore,
# by leaving it empty we can ensure that the container and binary shipped on it will have the same platform.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${GOARCH} go build -a -o lcmjob cmd/lifecycle-job/main.go

# Download vendor tooling for firmware management
FROM debian:bookworm-20240408-slim as tooling
WORKDIR /temp
RUN apt update && apt install curl -y
# Lenovo OneCli
RUN curl -sO https://download.lenovo.com/servers/mig/2024/02/23/59428/lnvgy_utl_lxce_onecli01t-4.4.0_linux_x86-64.tgz
RUN tar -xzf lnvgy_utl_lxce_onecli01t-4.4.0_linux_x86-64.tgz && rm lnvgy_utl_lxce_onecli01t-4.4.0_linux_x86-64.tgz

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM debian:bookworm-20240408-slim
WORKDIR /
COPY --from=builder /workspace/lcmjob .
COPY --from=tooling /temp /lenovo
USER 65532:65532

ENTRYPOINT ["/lcmjob"]