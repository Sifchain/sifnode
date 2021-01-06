#!/bin/bash

set -x
set -e

. $(dirname $0)/vagrantenv.sh
. ${BASEDIR}/test/integration/shell_utilities.sh

# Rebuild sifchain, but this time don't use validators

ADD_VALIDATOR_TO_WHITELIST= bash ${BASEDIR}/test/integration/setup_sifchain.sh && . $envexportfile

python3 $TEST_INTEGRATION_DIR/no_whitelisted_validators.py $NETDEF_JSON
