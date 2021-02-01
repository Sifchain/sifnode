#!/bin/bash

set -e

. $(dirname $0)/vagrantenv.sh
. ${TEST_INTEGRATION_DIR}/shell_utilities.sh

set -x

pkill -9 ebrelayer || true

mkdir -p $datadir/logs
set_persistant_env_var EBRELAYER_LOG $datadir/logs/ebrelayer.$(filenamedate).log $envexportfile

nohup $TEST_INTEGRATION_DIR/sifchain_run_ebrelayer.sh < /dev/null > $EBRELAYER_LOG 2>&1 &
set_persistant_env_var EBRELAYER_PID $! $envexportfile

( tail -n +1 -F $EBRELAYER_LOG & ) | grep -m 1 "Started Ethereum websocket with provider"