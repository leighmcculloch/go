#! /usr/bin/env bash

find . -name 'go.mod' -print0 | xargs -0 -I {} go mod edit \
  -dropreplace="github.com/stellar-modules/go/sdk" \
  -dropreplace="github.com/stellar-modules/go/exp" \
  -dropreplace="github.com/stellar-modules/go/services/internal" \
  {}
