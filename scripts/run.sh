#!/usr/bin/env bash

killall sifnoded
export DAEMON_HOME=$HOME/.sifnoded
export DAEMON_NAME=sifnoded
export DAEMON_ALLOW_DOWNLOAD_BINARIES=true
mkdir -p $DAEMON_HOME/cosmovisor/genesis/bin
cp $GOPATH/bin/sifnoded $DAEMON_HOME/cosmovisor/genesis/bin/

mkdir -p $DAEMON_HOME/cosmovisor/upgrades/

#sifnoded rest-server &
cosmovisor start
#--x-crisis-skip-assert-invariants --inv-check-period 600