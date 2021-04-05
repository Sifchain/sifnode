#!/bin/bash

# this script can be run against any chain.  All the tests create their own addresses
# and don't rely on preexisting state.

set -e

. $(dirname $0)/vagrantenv.sh
. $(dirname $0)/shell_utilities.sh

logecho $0 starting

loglevel=${LOG_LEVEL:-INFO}

python3 -m pytest -olog_level=$loglevel -v -olog_file=/tmp/log.txt -v \
  ${TEST_INTEGRATION_PY_DIR}/test_eth_transfers.py \
  ${TEST_INTEGRATION_PY_DIR}/test_rowan_transfers.py \
