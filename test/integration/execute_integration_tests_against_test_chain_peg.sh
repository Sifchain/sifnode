#!/bin/bash

# this script can be run only against a test chain.  It relies on tight control over
# ganache and knows where ganache stores its data.

set -ex

. $(dirname $0)/vagrantenv.sh
. $(dirname $0)/shell_utilities.sh

loglevel=${LOG_LEVEL:-INFO}

logecho $0 starting

env | sort

python3 -m pytest -olog_level=$loglevel -v -olog_file=/tmp/log.txt -v \
  ${TEST_INTEGRATION_PY_DIR}/test_liquidity_pools.py \
  ${TEST_INTEGRATION_PY_DIR}/test_new_currency_transfers.py \
  ${TEST_INTEGRATION_PY_DIR}/test_peggy_fees.py \
  ${TEST_INTEGRATION_PY_DIR}/test_random_currency_roundtrip.py \
  ${TEST_INTEGRATION_PY_DIR}/test_rollback_chain.py \

# run replay tests after other tests since they interact badly with replaydb

python3 -m pytest -olog_level=$loglevel -v -olog_file=/tmp/log.txt -v \
  ${TEST_INTEGRATION_PY_DIR}/test_ebrelayer_replay.py \
  ${TEST_INTEGRATION_PY_DIR}/test_ebrelayer_restart.py \

echo these tests do not run yet: \
  ${TEST_INTEGRATION_PY_DIR}/test_bulk_transfers.py \
