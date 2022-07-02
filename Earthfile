VERSION 0.6


COPY_METADATA:
    COMMAND
    COPY ".git" ".git"
    COPY "./ci" "./ci"
    COPY "./VERSION" "./VERSION"


clean-git-history-checking:
    FROM rust
    RUN cargo install clean_git_history
    DO +COPY_METADATA
    ARG from_reference="origin/HEAD"
    RUN ./ci/clean-git-history-checking.sh --from-reference "${from_reference}"


conventional-commits-linting:
    FROM rust
    RUN cargo install conventional_commits_linter
    DO +COPY_METADATA
    ARG from_reference="origin/HEAD"
    RUN ./ci/conventional-commits-linting.sh --from-reference "${from_reference}"


conventional-commits-next-version-checking:
    FROM rust
    RUN cargo install conventional_commits_next_version
    DO +COPY_METADATA
    RUN ./ci/conventional-commits-next-version-checking.sh


COPY_SOURCECODE:
    COMMAND
    COPY "go.mod" "go.mod"
    COPY "go.sum" "go.sum"
    COPY "./ci" "./ci"
    COPY "./src" "./src"


SAVE_OUTPUT:
    COMMAND
    SAVE ARTIFACT "git-manager" AS LOCAL "git-manager"
    SAVE ARTIFACT "go.sum" AS LOCAL "go.sum"


golang-base:
    FROM golang:1.18
    WORKDIR /tmp/git-manager
    ENV GOPROXY=direct
    ENV CGO_ENABLED=0
    ENV GOOS=linux
    ENV GOARCH=amd64


check-formatting:
    FROM +golang-base
    DO +COPY_SOURCECODE
    RUN ./ci/check-formatting.sh


fix-formatting:
    FROM +golang-base
    DO +COPY_SOURCECODE
    RUN ./ci/fix-formatting.sh
    SAVE ARTIFACT "./src" AS LOCAL "./src"


linting:
    FROM +golang-base
    RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.47.2
    DO +COPY_SOURCECODE
    RUN ./ci/linting.sh


check-module-tidying:
    FROM +golang-base
    DO +COPY_SOURCECODE
    RUN ./ci/check-module-tidying.sh


fix-module-tidying:
    FROM +golang-base
    DO +COPY_SOURCECODE
    RUN ./ci/fix-module-tidying.sh
    SAVE ARTIFACT "go.mod" AS LOCAL "go.mod"


compiling-linux-amd64:
    FROM +golang-base
    DO +COPY_SOURCECODE
    RUN ./ci/compiling.sh
    DO +SAVE_OUTPUT


compiling-darwin-amd64:
    FROM +golang-base
    ENV GOOS=darwin
    DO +COPY_SOURCECODE
    RUN ./ci/compiling.sh
    DO +SAVE_OUTPUT


releasing:
    FROM rust
	# Install release description generator.
	RUN cargo install git-cliff
	# Install GitlabCI cli releasing tool.
	RUN curl --location --output /usr/local/bin/release-cli "https://release-cli-downloads.s3.amazonaws.com/latest/release-cli-linux-amd64"
	RUN chmod +x /usr/local/bin/release-cli
    DO +COPY_METADATA
    ARG server_url
    ARG job_token
    ARG project_id
    ARG reference
    RUN ./ci/releasing.sh "${server_url}" "${job_token}" "${project_id}" "${reference}"
