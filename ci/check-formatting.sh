#!/usr/bin/env sh

set -o errexit
set -o xtrace

output=$(gofmt -d -l "./src")
echo "${output}"
test -z "${output}"
