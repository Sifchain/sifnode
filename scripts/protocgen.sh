#!/usr/bin/env bash

set -eo pipefail

protoc_gen_gocosmos() {
  if ! grep "github.com/gogo/protobuf => github.com/regen-network/protobuf" go.mod &>/dev/null ; then
    echo -e "\tPlease run this command from somewhere inside the cosmos-sdk folder."
    return 1
  fi

  go get github.com/regen-network/cosmos-proto/protoc-gen-gocosmos@latest 2>/dev/null
}

protoc_gen_gocosmos

buf generate

# command to generate docs using protoc-gen-doc
# buf protoc \
#   -I "proto" \
#   -I "third_party/proto" \
#   --doc_out=./docs/core \
#   --doc_opt=./docs/protodoc-markdown.tmpl,proto-docs.md \
#   $(find "$(pwd)/proto" -maxdepth 5 -name '*.proto')
# go mod tidy

# move proto files to the right places
cp -r github.com/Sifchain/sifnode/* ./
rm -rf github.com
