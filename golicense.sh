#! /bin/bash
set -e

printf "Running go license...\n"
command -v go-licenses >/dev/null 2>&1 || (
  dir=$(mktemp -d)
  pushd $dir
  go mod init golicense
  go get github.com/google/go-licenses
  popd
)

go-licenses csv ./... > go.licenses
