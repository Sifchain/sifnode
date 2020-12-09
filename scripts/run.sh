#!/usr/bin/env bash

killall sifnoded sifnodecli

sifnoded start
#>> sifnode.log 2>&1  &
sifnodecli rest-server &
