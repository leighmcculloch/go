#! /usr/bin/env bash

function modcreate() {
  local path="$1"
  local pathmod="$path/go.mod"

  # copy the root go.mod so we preserve consistency of existing min versions
  cp go.mod $pathmod
  go mod edit -module="github.com/stellar/go/$path" "$pathmod"

  # add require and replace statements for the new modules in the repo that are
  # dependencies of this module
  local replacements=(${@:2})
  local root="$(realpath .)"
  for replacement in "${replacements[@]}"
  do
    if [ -z "$replacement" ]
    then
      continue
    fi
    if [ "$replacement" == "." ]
    then
      continue
    fi
    local relpath="$(realpath --relative-to="$path" "$root/$replacement")"
    go mod edit -require="github.com/stellar/go/$replacement@v0.0.0-00010101000000-000000000000" "$pathmod"
    go mod edit -replace="github.com/stellar/go/$replacement"="$relpath" "$pathmod"
  done

  # this code below is commented out because the reference back is causing
  # issues resolving other dependencies. it's unclear why.

  ## add a dependency on the root module to prevent users of the new modules
  ## from accidentally importing an old version of the root that contains the
  ## same import paths, causing ambiguous import path errors
  #local rootrelpath="$(realpath --relative-to="$path" "$root")"
  #go mod edit -require="github.com/stellar/go@v0.0.0-00010101000000-000000000000" "$pathmod"
  #go mod edit -replace="github.com/stellar/go"="$rootrelpath/" "$pathmod"

  ## add a go file with blank import path to teach the go toolchain to keep the
  ## require on the root repo that we added in the lines immediately above
  #pushd $path >> /dev/null
  #local package="$(go list -f '{{.Name}}' || echo ${PWD##*/})"
  #local modroot="package $package\nimport _ \"github.com/stellar/go\""
  #echo -e "$modroot" | gofmt > stellargo.go
  #popd >> /dev/null

  pushd $path >> /dev/null
  go mod tidy
  popd >> /dev/null
}

function modtidy() {
  find . -name 'go.mod' -printf "%h\0" | xargs -0 -I {} bash -c "cd {}; pwd; go mod tidy; go list -m all > go.list"
}

modcreate 'sdk'

modcreate 'exp' sdk

modcreate 'services/internal' sdk

modcreate 'services/bridge' sdk services/internal
modcreate 'services/compliance' sdk services/internal
modcreate 'services/federation' sdk
modcreate 'services/friendbot' sdk
modcreate 'services/horizon' sdk exp
modcreate 'services/keystore' sdk
modcreate 'services/ticker' sdk

modcreate 'tools/archive-reader' sdk exp
modcreate 'tools/horizon-cmp' sdk
modcreate 'tools/horizon-verify' sdk
modcreate 'tools/stellar-archivist' sdk
modcreate 'tools/stellar-hd-wallet' sdk exp
modcreate 'tools/stellar-sign' sdk
modcreate 'tools/stellar-vanity-gen' sdk

# this code below is commented out because of the reason above
#echo "package stellargo" | gofmt > "root.go"

modtidy
modtidy

megalist=$(find . -name 'go.list' -print0 | tac -s $'\0' | xargs -0 -I {} cat {} | LC_ALL=C sort -u | grep -v -E 'github\.com/stellar/go([ /]|$)')
masterlist=$(git show master:go.list | grep -v -E 'github\.com/stellar/go([ /]|$)')
colordiff -u <(echo -n "$masterlist") <(echo -n "$megalist")
