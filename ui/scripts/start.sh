#!/bin/bash

set -e

yarn concurrently -k -r -s first "./scripts/_eth-revert.sh" "yarn wait-on tcp:localhost:7545 && ./scripts/_sif-revert.sh" "yarn wait-on http-get://localhost:1317/node_info tcp:localhost:7545 && ./scripts/_peggy-revert.sh"