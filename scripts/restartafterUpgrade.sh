#!/usr/bin/env bash

cp $GOPATH/src/new/sifnoded $GOPATH/bin/
cosmovisor start >> sifnode.log 2>&1  &