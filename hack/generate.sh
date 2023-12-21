#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_DIR="$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
export TERM="xterm-256color"

bold="$(tput bold)"
blue="$(tput setaf 4)"
normal="$(tput sgr0)"

function qualify-gvs() {
  APIS_PKG="$1"
  GROUPS_WITH_VERSIONS="$2"
  join_char=""
  res=""

  for GVs in ${GROUPS_WITH_VERSIONS}; do
    IFS=: read -r G Vs <<<"${GVs}"

    for V in ${Vs//,/ }; do
      res="$res$join_char$APIS_PKG/$G/$V"
      join_char=","
    done
  done

  echo "$res"
}

function qualify-gs() {
  APIS_PKG="$1"
  unset GROUPS
  IFS=' ' read -ra GROUPS <<< "$2"
  join_char=""
  res=""

  for G in "${GROUPS[@]}"; do
    res="$res$join_char$APIS_PKG/$G"
    join_char=","
  done

  echo "$res"
}

function generate() {
  package="$1"
  (
  cd "$VIRTUAL_GOPATH/src"
  export PATH="$PATH:$(dirname "$PROTOC_GEN_GOGO")"
  echo "Generating ${blue}$package${normal}"
  protoc \
    --proto_path "./github.com/ironcore-dev/lifecycle-manager/$package" \
    --proto_path "$VIRTUAL_GOPATH/src" \
    --gogo_out=plugins=grpc:"$VIRTUAL_GOPATH/src" \
    "./github.com/ironcore-dev/lifecycle-manager/$package/api.proto"
  )
}

VGOPATH="$VGOPATH"
DEEPCOPY_GEN="$DEEPCOPY_GEN"
GO_TO_PROTOBUF="$GO_TO_PROTOBUF"
PROTOC_GEN_GOGO="$PROTOC_GEN_GOGO"

VIRTUAL_GOPATH="$(mktemp -d)"
trap 'rm -rf "$VIRTUAL_GOPATH"' EXIT

# Setup virtual GOPATH so the codegen tools work as expected.
(cd "$SCRIPT_DIR/.."; go mod download && "$VGOPATH" -o "$VIRTUAL_GOPATH")

export GOROOT="${GOROOT:-"$(go env GOROOT)"}"
export GOPATH="$VIRTUAL_GOPATH"
export GO111MODULE=off

CLIENT_GROUPS="lifecycle"
CLIENT_VERSION_GROUPS="lifecycle:v1alpha1"
ALL_VERSION_GROUPS="$CLIENT_VERSION_GROUPS"

echo "${bold}Public types${normal}"

echo "Generating ${blue}deepcopy${normal}"
"$DEEPCOPY_GEN" \
  --output-base "$GOPATH/src" \
  --go-header-file "$SCRIPT_DIR/boilerplate.go.txt" \
  --input-dirs "$(qualify-gvs "github.com/ironcore-dev/lifecycle-manager/api" "$ALL_VERSION_GROUPS")" \
  -O zz_generated.deepcopy

echo "Generating ${blue}protobuf${normal}"
"$GO_TO_PROTOBUF" \
  --output-base "$GOPATH/src" \
  --go-header-file "$SCRIPT_DIR/boilerplate.go.txt" \
  --packages "$(qualify-gvs "github.com/ironcore-dev/lifecycle-manager/api" "$ALL_VERSION_GROUPS")" \
  --apimachinery-packages "-k8s.io/apimachinery/pkg/runtime/schema,-k8s.io/apimachinery/pkg/apis/meta/v1,-k8s.io/api/core/v1"

generate "lcmi/api/common/v1alpha1"
generate "lcmi/api/machine/v1alpha1"
generate "lcmi/api/machinetype/v1alpha1"
