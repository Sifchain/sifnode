#!/bin/bash

#
# Sifnode entrypoint.
#

set -x

. $(dirname $0)/vagrantenv.sh
. $TEST_INTEGRATION_DIR/shell_utilities.sh

set_persistant_env_var GANACHE_LOG $datadir/logs/ganache.$(filenamedate).txt $envexportfile
mkdir -p $(dirname $GANACHE_LOG)

pkill -f -9 ganache-cli || true
while nc -z localhost 7545; do
  sleep 1
done

nohup tmux new-session -d -s my_session "ganache-cli -h 0.0.0.0 --mnemonic 'candy maple cake sugar pudding cream honey rich smooth crumble sweet treat' --networkId '5777' --port '7545' --db ${GANACHE_DB_DIR} > $GANACHE_LOG 2>&1"

sleep 5

while ! nc -z localhost 7545; do
  sleep 5
done
