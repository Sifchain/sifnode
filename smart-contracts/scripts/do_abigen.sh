#!/bin/bash

set -e

jsonfile=$1
shift
outputbase=$1
shift

jsondir=$(dirname $jsonfile)
pkg=$(basename "$jsonfile" .json)
outputdir=${outputbase}/$jsondir
mkdir -p $outputdir
# jq .abi < artifacts/contracts/BridgeBank/BridgeBank.sol/BridgeBank.json | abigen --abi - --pkg Foo
jq .abi < $jsonfile | abigen --abi - --pkg $pkg --type $pkg --out ${outputdir}/$pkg.go
jq .abi < $jsonfile | abigen --abi - --pkg $pkg --type $pkg --out ${jsondir}/$pkg.go
