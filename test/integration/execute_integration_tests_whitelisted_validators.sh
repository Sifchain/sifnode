#!/bin/bash

set -e

. $(dirname $0)/vagrantenv.sh
. ${BASEDIR}/test/integration/shell_utilities.sh

loglevel=${LOG_LEVEL:-INFO}

# Rebuild sifchain, but this time don't use validators

ADD_VALIDATOR_TO_WHITELIST= bash ${BASEDIR}/test/integration/setup_sifchain.sh && . $envexportfile

# TODO we should get this from a script, not hardcoded
operator_address=0xf17f52151EbEF6C7334FAD080c5704D77216b732

ETHEREUM_ADDRESS=$operator_address python3 -m pytest -olog_level=$loglevel -olog_cli=true -v -olog_file=/tmp/log.txt -v \
  $TEST_INTEGRATION_PY_DIR/no_whitelisted_validators.py

# rebuild again with validators so the chain is useful for other things
ADD_VALIDATOR_TO_WHITELIST=true bash ${BASEDIR}/test/integration/setup_sifchain.sh && . $envexportfile
