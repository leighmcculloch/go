#! /bin/bash
set -e

./goall.sh go mod tidy
./goall.sh git diff --exit-code -- go.mod || (echo "Go file go.mod is dirty, update the file with './goall.sh go mod tidy' locally." && exit 1)
./goall.sh git diff --exit-code -- go.sum || (echo "Go file go.sum is dirty, update the file with './goall.sh go mod tidy' locally." && exit 1)
./goall.sh bash -c 'diff -u go.list <(go list -m all)' || (echo "Go dependencies have changed, update the go.list file with './goall.sh bash -c \"go list -m all > go.list\"' locally." && exit 1)
./goall.sh go mod verify || (echo "One or more Go dependencies failed verification. Either a version is no longer available, or the author or someone else has modified the version so it no longer points to the same code." && exit 1)
