#!/usr/bin/env sh

set -o errexit
set -o xtrace

gofmt -w -s "./src"
