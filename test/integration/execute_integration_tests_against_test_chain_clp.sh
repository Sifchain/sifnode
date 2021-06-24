#!/bin/bash

# this script can be run only against a test chain.  It relies on tight control over
# ganache and knows where ganache stores its data.

set -e

. $(dirname $0)/vagrantenv.sh
. $(dirname $0)/shell_utilities.sh

loglevel=${LOG_LEVEL:-INFO}

logecho $0 starting

python3 -m pytest -olog_level=$loglevel -v -olog_file=/tmp/log.txt -v \
  ${TEST_INTEGRATION_PY_DIR}/test_dispensation_offlinetxn.py \
  ${TEST_INTEGRATION_PY_DIR}/test_dispensation_onlinetxn.py \
  ${TEST_INTEGRATION_PY_DIR}/test_claims.py \
  ${TEST_INTEGRATION_PY_DIR}/test_dispensation_volume9.py \
  ${TEST_INTEGRATION_PY_DIR}/test_dispensation_volume12.py \
  ${TEST_INTEGRATION_PY_DIR}/test_liquidity_pools.py \
