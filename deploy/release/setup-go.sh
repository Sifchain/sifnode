#!/bin/bash -eux

GO_LINUX_PACKAGE_URL="https://dl.google.com/go/go1.15.1.linux-amd64.tar.gz"

wget --progress=dot:mega ${GO_LINUX_PACKAGE_URL} -O go-linux.tar.gz
tar -zxf go-linux.tar.gz
mv go /usr/local/
mkdir -p /go/bin /go/src /go/pkg

export GO_HOME=/usr/local/go
export GOPATH=/go
export PATH=${GOPATH}/bin:${GO_HOME}/bin/:$PATH
