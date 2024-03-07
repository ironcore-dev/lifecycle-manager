#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_DIR="$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
export TERM="xterm-256color"

bold="$(tput bold)"
blue="$(tput setaf 4)"
normal="$(tput sgr0)"

trap 'popd > /dev/null' EXIT

pushd "$SCRIPT_DIR/.." > /dev/null
ROOT=$(pwd)
popd > /dev/null

pushd "$SCRIPT_DIR/../lcmi" > /dev/null
export PATH=$PATH:$ROOT/bin

echo "Generating ${blue}proto${normal}"
buf generate
