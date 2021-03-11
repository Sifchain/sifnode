#!/bin/bash

set -ev

. $(dirname $0)/vagrantenv.sh
. ${TEST_INTEGRATION_DIR}/shell_utilities.sh

set -x

pkill -9 ebrelayer || true

mkdir -p $datadir/logs
set_persistant_env_var EBRELAYER_LOG $datadir/logs/ebrelayer.$(filenamedate).log $envexportfile

nohup $TEST_INTEGRATION_DIR/sifchain_run_ebrelayer.sh < /dev/null > $EBRELAYER_LOG 2>&1 &
set_persistant_env_var EBRELAYER_PID $! $envexportfile

# This doesn't work either (from python), although it worked from bash on Macbook.  Commenting it out for now as we really don't need it.
# fail and add timeout to the check.
#timeout 30s grep -m 1 'Started Ethereum websocket with provider' <(tail -n +1 -f $EBRELAYER_LOG) || exit 1
# This doesn't work on macbook for some reason
#( tail -n +1 -F $EBRELAYER_LOG & ) | grep -m 1 "Started Ethereum websocket with provider"