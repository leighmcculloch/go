#! /usr/bin/env bash

function pkgmove() {
  local path="$1"
  local pathdest="$2"
  mv "$path" "$pathdest"
  find . -type f -not -path '*/.git/*' -print0 | xargs -0 sed -i "s#github\.com/stellar/go/$path#github\.com/stellar/go/$pathdest#"
}

mkdir sdk
pkgmove 'address' 'sdk/address'
pkgmove 'amount' 'sdk/amount'
pkgmove 'build' 'sdk/build'
pkgmove 'clients' 'sdk/clients'
pkgmove 'crc16' 'sdk/crc16'
pkgmove 'handlers' 'sdk/handler'
pkgmove 'hash' 'sdk/hash'
pkgmove 'keypair' 'sdk/keypair'
pkgmove 'meta' 'sdk/meta'
pkgmove 'network' 'sdk/network'
pkgmove 'price' 'sdk/price'
pkgmove 'protocols' 'sdk/protocols'
pkgmove 'strkey' 'sdk/strkey'
pkgmove 'support' 'sdk/support'
pkgmove 'txnbuild' 'sdk/txnbuild'
pkgmove 'xdr' 'sdk/xdr'
mv 'main_test.go' 'sdk/'
sed -i '1 s/package stellargo/package sdk/' 'sdk/main_test.go'

go mod tidy
