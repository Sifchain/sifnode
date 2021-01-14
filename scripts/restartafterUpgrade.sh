#!/usr/bin/env bash

cp $GOPATH/src/new/sifnodecli $GOPATH/bin/
cosmovisor start >> sifnode.log 2>&1  &