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

parallelizr "sifnoded start" "sifnodecli rest-server  --unsafe-cors --trace"