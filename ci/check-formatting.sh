#!/usr/bin/env sh

set -o errexit
set -o xtrace

output=$(gofmt -l ".")
echo "${output}"
test -z "${output}"
