#!/bin/bash

set -x
set -e

. $(dirname $0)/vagrantenv.sh
. $(dirname $0)/shell_utilities.sh

logecho $0 starting

python3 $TEST_INTEGRATION_DIR/initial_test_balances.py $NETDEF_JSON
logecho $0 completed $TEST_INTEGRATION_DIR/initial_test_balances.py
sleep 15
python3 $TEST_INTEGRATION_DIR/peggy-basic-test-docker.py $NETDEF_JSON
python3 $TEST_INTEGRATION_DIR/peggy-e2e-test.py $NETDEF_JSON
python3 $TEST_INTEGRATION_DIR/test_chain_rollback.py $NETDEF_JSON

# save sifchain transaction data; later whitelist testing will intentionally break transaction queries
bash $TEST_INTEGRATION_DIR/sifchain_logs.sh
