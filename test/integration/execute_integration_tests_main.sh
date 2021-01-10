#!/bin/bash

set -x
set -e

. $(dirname $0)/vagrantenv.sh
. $(dirname $0)/shell_utilities.sh

logecho $0 starting

# TODO we should get these from a script, not hardcoded
operator_address=0xf17f52151EbEF6C7334FAD080c5704D77216b732

ETHEREUM_ADDRESS=$operator_address python3 -m pytest -o=log_cli=true -olog_level=DEBUG $TEST_INTEGRATION_DIR/test_new_account.py

python3 $TEST_INTEGRATION_DIR/initial_test_balances.py $NETDEF_JSON
python3 $TEST_INTEGRATION_DIR/peggy-basic-test-docker.py $NETDEF_JSON
python3 $TEST_INTEGRATION_DIR/peggy-e2e-test.py $NETDEF_JSON
python3 $TEST_INTEGRATION_DIR/test_chain_rollback.py $NETDEF_JSON

# save sifchain transaction data; later whitelist testing will intentionally break transaction queries
bash $TEST_INTEGRATION_DIR/sifchain_logs.sh
