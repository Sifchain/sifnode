#!/bin/bash

killall sifnoded sifnodecli ebrelayer ganache-cli

yarn concurrently -k -r -s first \
  "yarn chain:eth:revert" \
  "yarn wait-on tcp:localhost:7545 && yarn chain:sif:revert" \
  "yarn wait-on http-get://localhost:1317/node_info tcp:localhost:7545 && yarn chain:peggy:revert"