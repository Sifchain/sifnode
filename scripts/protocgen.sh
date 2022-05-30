#!/usr/bin/env bash

set -eo pipefail

go get github.com/regen-network/cosmos-proto/protoc-gen-gocosmos@latest 2>/dev/null
buf generate

# move proto files to the right places
cp -r github.com/Sifchain/sifnode/* ./
rm -rf github.com google.golang.org
