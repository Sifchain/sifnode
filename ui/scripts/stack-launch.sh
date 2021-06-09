#!/bin/bash


# This will run our stack from genesis
# All setup scripts will be run and once it is ready our app will be compiled and built
# You probably only need to use this script to set a new snapshot archive

killall sifnoded sifnodecli ebrelayer ganache-cli
sleep 5

./scripts/_sif-build.sh

yarn concurrently -k -r -s first \
  "./scripts/_eth.sh" \
  "./scripts/_sif.sh" \
  "yarn wait-on http-get://localhost:1317/node_info && ./scripts/_migrate.sh && ./scripts/_peggy.sh"
