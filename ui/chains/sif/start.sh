#!/bin/bash

parallelizr() {
  for cmd in "$@"; do {
    $cmd & pid=$!
    PID_LIST+=" $pid";
  } done

  trap "kill -9 $PID_LIST" SIGINT

  wait $PID_LIST
}

echo "Starting test chain"

parallelizr "sifnoded start --minimum-gas-prices 0.5rowan --rpc.laddr tcp://0.0.0.0:26657" "sifnodecli rest-server --laddr tcp://0.0.0.0:1317 --unsafe-cors --trace"