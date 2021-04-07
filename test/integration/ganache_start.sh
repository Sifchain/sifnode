#!/bin/bash

# $0 block-delay-time [default is none, always generate a block]
#   block-delay-time is passed to ganache as -b block-delay-time

block_delay=$1
if [ ! -z "$block_delay" ]
then
    block_delay="-b $block_delay"
fi

. $(dirname $0)/vagrantenv.sh
. $TEST_INTEGRATION_DIR/shell_utilities.sh

set_persistant_env_var GANACHE_LOG $datadir/logs/ganache.$(filenamedate).txt $envexportfile
set_persistant_env_var GANACHE_KEYS_JSON $datadir/ganachekeys.json $envexportfile

rm -f $GANACHE_KEYS_JSON

mkdir -p $(dirname $GANACHE_LOG)

pkill -9 -f ganache-cli || true
while nc -z localhost 7545; do
  sleep 1
done


# ganache really hates running in the background.  Put it in a tmux session to keep all its input code happy.
# If you don't do this, ganache-cli will just exit.
nohup tmux new-session -d -s my_session "ganache-cli ${block_delay} -h 0.0.0.0 --mnemonic 'candy maple cake sugar pudding cream honey rich smooth crumble sweet treat' --networkId '5777' --port '7545' --db ${GANACHE_DB_DIR} --account_keys_path $GANACHE_KEYS_JSON > $GANACHE_LOG 2>&1"

# wait for ganache to come up
sleep 5

while ! nc -z localhost 7545; do
  sleep 5
done
while [ ! -f $GANACHE_KEYS_JSON ]; do
  sleep 1
done
