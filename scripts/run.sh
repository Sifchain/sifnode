#!/usr/bin/env bash

killall sifnoded sifnodecli

sifnoded start &
sifnodecli rest-server &
