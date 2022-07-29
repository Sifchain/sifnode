#!/usr/bin/env bash

killall sifnoded

cd ../..
make install
sifnoded start --trace
