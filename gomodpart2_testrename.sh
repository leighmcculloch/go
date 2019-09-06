#! /usr/bin/env bash

find . -name '*.mod' -print0 | xargs -0 sed -i -E 's#github\.com/stellar/go$#github.com/stellar-modules/go#'

find . -name '*.go' -print0 | xargs -0 sed -i -E 's#github\.com/stellar/go([" /])#github.com/stellar-modules/go\1#'
find . -name '*.mod' -print0 | xargs -0 sed -i -E 's#github\.com/stellar/go([" /])#github.com/stellar-modules/go\1#'
find . -name '*.list' -print0 | xargs -0 sed -i -E 's#github\.com/stellar/go([" /])#github.com/stellar-modules/go\1#'
