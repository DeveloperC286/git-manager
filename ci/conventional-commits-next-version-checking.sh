#!/usr/bin/env sh

set -o errexit
set -o xtrace

# Get current version.
current_version=$(cat "VERSION")
# Get latest tag.
latest_tag=$(git tag --sort=-committerdate | head -1)
# Check current vs expected.
"${CARGO_HOME}/bin/conventional_commits_next_version" --calculation-mode "Batch" --from-reference "${latest_tag}" --from-version "${latest_tag}" --current-version "${current_version}"
