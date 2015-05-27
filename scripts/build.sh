#!/bin/env sh

# any commands that fail cause the script to exit
set -e

# Get the git commit
GIT_COMMIT=$(git rev-parse HEAD)
GIT_DIRTY=$(test -n "`git status --porcelain`" && echo "+CHANGES" || true)

go build -v
cd cmd/plumb
go build -v -ldflags "-X main.GitCommit ${GIT_COMMIT}${GIT_DIRTY}"
cd ../..
mv cmd/plumb/plumb .
