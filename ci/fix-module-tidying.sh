#!/usr/bin/env sh

set -o errexit
set -o xtrace

go mod tidy -v
