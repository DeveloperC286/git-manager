VERSION 0.6


clean-git-history-checking:
    FROM rust
    RUN cargo install clean_git_history
	COPY ".git" "."
	ARG from="origin/HEAD"
	RUN /usr/local/cargo/bin/clean_git_history --from-reference "${from}"


conventional-commits-linting:
    FROM rust
    RUN cargo install conventional_commits_linter
    COPY ".git" "."
    ARG from="origin/HEAD"
    RUN /usr/local/cargo/bin/conventional_commits_linter --from-reference "${from}" --allow-angular-type-only


compiling:
    FROM golang:1.18
	WORKDIR /tmp/git-manager
	ENV GOPROXY direct
	ENV CGO_ENABLED 0
	ENV GOOS darwin
	ENV GOARCH amd64
    COPY "go.mod" "go.mod"
    COPY "go.sum" "go.sum"
    # Need latest verison otherwise get 'go:linkname must refer to declared function or variable'
	# See: https://stackoverflow.com/questions/71507321/go-1-18-build-error-on-mac-unix-syscall-darwin-1-13-go253-golinkname-mus
	RUN go get -u golang.org/x/sys
    # Copy in source last so other steps are cached.
    COPY "./src" "./src"
    RUN go build -o git-manager "./src/"
    SAVE ARTIFACT "git-manager" AS LOCAL "git-manager"
