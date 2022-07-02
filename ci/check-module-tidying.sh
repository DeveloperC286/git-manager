#!/usr/bin/env sh

set -o errexit
set -o xtrace

cp "go.mod" "go.mod.orig"
go mod tidy -v
cmp "go.mod" "go.mod.orig"
exit_code=$?
# Cleanup.
mv "go.mod.orig" "go.mod"
exit ${exit_code}
