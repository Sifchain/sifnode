#!/bin/bash

set -x
set -e

. $(dirname $0)/vagrantenv.sh
. ${TEST_INTEGRATION_DIR}/shell_utilities.sh

pkill sifnoded || true
pkill ebrelayer || true

sleep 1

#
# Remove prior generations Config
#
if [ -d $NETWORKDIR ]
then
  # $NETWORKDIR has many directories without write permission, so change those
  # before deleting it.
  find $NETWORKDIR -type d | xargs chmod +w
  rm -rf $NETWORKDIR && mkdir $NETWORKDIR
fi
mkdir -p $NETWORKDIR
sifgen network create localnet 1 $NETWORKDIR 192.168.1.2 $NETWORKDIR/network-definition.yml --keyring-backend file  --mint-amount 999999000000000000000000000rowan,1370000000000000000ibc/FEEDFACEFEEDFACEFEEDFACEFEEDFACEFEEDFACEFEEDFACEFEEDFACEFEEDFACE

set_persistant_env_var NETDEF $NETWORKDIR/network-definition.yml $envexportfile
set_persistant_env_var NETDEF_JSON $datadir/netdef.json $envexportfile
cat $NETDEF | to_json > $NETDEF_JSON

set_persistant_env_var MONIKER $(cat $NETDEF_JSON | jq -r '.[0].moniker') $envexportfile
set_persistant_env_var VALIDATOR1_PASSWORD $(cat $NETDEF_JSON | jq -r '.[0].password') $envexportfile
set_persistant_env_var VALIDATOR1_ADDR $(cat $NETDEF_JSON | jq -r '.[0].address') $envexportfile
set_persistant_env_var MNEMONIC "$(cat $NETDEF_JSON | jq -r '.[0].mnemonic')" $envexportfile
set_persistant_env_var CHAINDIR $NETWORKDIR/validators/$CHAINNET/$MONIKER $envexportfile
set_persistant_env_var SIFNODED_LOG $datadir/logs/sifnoded.log $envexportfile

mkdir -p $datadir/logs
nohup $TEST_INTEGRATION_DIR/sifchain_start_daemon.sh < /dev/null > $SIFNODED_LOG 2>&1 &
set_persistant_env_var SIFNODED_PID $! $envexportfile
nohup sifnoded rest-server --laddr tcp://0.0.0.0:1317 < /dev/null > $datadir/logs/restserver.log 2>&1 &
set_persistant_env_var REST_SERVER_PID $! $envexportfile
bash $TEST_INTEGRATION_DIR/sifchain_start_ebrelayer.sh
