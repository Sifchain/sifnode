#!/bin/bash

# this script can be run only against a test chain.  It relies on tight control over
# ganache and knows where ganache stores its data.

set -e

. $(dirname $0)/vagrantenv.sh
. $(dirname $0)/shell_utilities.sh

loglevel=${LOG_LEVEL:-INFO}

logecho $0 starting

# TODO we should get this from a script, not hardcoded
operator_address=0xf17f52151EbEF6C7334FAD080c5704D77216b732

ETHEREUM_ADDRESS=$operator_address python3 -m pytest -olog_level=$loglevel -olog_cli=true -v -olog_file=/tmp/log.txt -v \
  ${TEST_INTEGRATION_PY_DIR}/test_rowan_transfers.py \
  ${TEST_INTEGRATION_PY_DIR}/test_rollback_chain.py
