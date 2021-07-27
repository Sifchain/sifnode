#!/usr/bin/env bash

killall sifnoded

sifnoded rest-server &
sifnoded start

