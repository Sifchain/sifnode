#!/bin/bash

# this script can be run against any chain.  All the tests create their own addresses
# and don't rely on preexisting state.

set -e

. $(dirname $0)/vagrantenv.sh
. $(dirname $0)/shell_utilities.sh

logecho $0 starting

loglevel=${LOG_LEVEL:-INFO}

# TODO we should get this from a script, not hardcoded
operator_address=0xf17f52151EbEF6C7334FAD080c5704D77216b732

ETHEREUM_ADDRESS=$operator_address python3 -m pytest -olog_level=$loglevel -olog_cli=true -v -olog_file=/tmp/log.txt -v \
  ${TEST_INTEGRATION_PY_DIR}/test_transfers.py \
  ${TEST_INTEGRATION_PY_DIR}/test_eth_transfers.py