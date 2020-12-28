#!/usr/bin/env bash

killall sifnoded sifnodecli

sifnodecli rest-server &
sifnoded start

