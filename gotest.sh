#! /bin/bash
set -e

if [ -z "$GOTESTCOVERPROFILE" ]
then
  ./goall.sh go test -race -covermode=atomic ./...
else
  ./goall.sh go test -race -covermode=atomic -coverprofile="$GOTESTCOVERPROFILE" ./...
fi
