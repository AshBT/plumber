#!/bin/bash

# any commands that fail cause the script to exit
set -e
if [ -z $1 ]; then
  echo "We expected at least one argument"
  exit 1
fi

if [[ $1 == "build" || $1 == "install" ]]; then
  # Get the git commit
  GIT_COMMIT=$(git describe --always)

  go $1 -v -ldflags "-X main.GitCommit ${GIT_COMMIT}"
else
  echo "We can only 'build' or 'install': got '$1' instead."
fi
