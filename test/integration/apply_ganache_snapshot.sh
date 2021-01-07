#!/bin/bash

# Takes a snapshot of the ganache data directory and prints
# the snapshot directory

set -ex

SNAPSHOT_DB_DIR=$1
shift

. $(dirname $0)/vagrantenv.sh
. $TEST_INTEGRATION_DIR/shell_utilities.sh

set_persistant_env_var GANACHE_DB_DIR $SNAPSHOT_DB_DIR $envexportfile

# save the previous ganache log
docker logs -t genesis_ganachecli_1 > ${datadir}/ganachelog.txt.$(filenamedate) 2>&1

logecho $0 restart ganache with snapshot

docker-compose --project-name genesis -f ${TEST_INTEGRATION_DIR}/docker-compose-ganache.yml down
docker-compose --project-name genesis -f ${TEST_INTEGRATION_DIR}/docker-compose-ganache.yml up -d --force-recreate

logecho $0 complete