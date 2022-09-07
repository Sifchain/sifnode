#!/usr/bin/env bash

set -x

killall sifnoded

cd ../..
make install
sifnoded start --trace
