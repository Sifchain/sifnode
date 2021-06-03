#!/bin/bash


# This will run our stack from genesis
# All setup scripts will be run and once it is ready our app will be compiled and built
# You probably only need to use this script to set a new snapshot archive

killall sifnoded sifnodecli ebrelayer ganache-cli
sleep 5

yarn concurrently -k -r -s first \
  "yarn chain:eth" \
  "yarn chain:sif" \
  "yarn wait-on http-get://localhost:1317/node_info && yarn chain:migrate && yarn chain:peggy"
