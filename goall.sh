#!/usr/bin/env bash

set -e

modules=()
readarray -d '' modules < <(find . -name 'go.mod' -printf "%h\0")

failures=()
for m in ${modules[@]}
do
  echo "# $m"
  pushd "$m" > /dev/null

  ec=0
  $@ || ec=$?
  if [ $ec -ne 0 ]
  then
    failures+=("$m exit code $ec")
  fi

  echo

  popd > /dev/null
done

if [ ${#failures[@]} -ne 0 ]
then
  echo FAILING MODULES:
  for failure in "${failures[@]}"
  do
    echo $failure
  done
  exit 1
fi
