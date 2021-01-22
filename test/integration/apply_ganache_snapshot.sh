#!/bin/bash

# Takes a snapshot of the ganache data directory and prints
# the snapshot directory

set -ex

SNAPSHOT_DB_DIR=$1
shift

. $(dirname $0)/vagrantenv.sh
. $TEST_INTEGRATION_DIR/shell_utilities.sh

set_persistant_env_var GANACHE_DB_DIR $SNAPSHOT_DB_DIR $envexportfile

logecho $0 restart ganache with snapshot

bash $TEST_INTEGRATION_DIR/ganache_start.sh

logecho $0 complete