#!/bin/bash
# must run from the root directory of the sifnode tree

set -e # exit on any failure

basedir=$(dirname $0)
pidfile=${basedir}/watcher_pids.txt

bash $basedir/stop_watchers.sh

docker logs -f integration_sifnode1_1 > ${datadir}/integrationlog.txt 2>&1 &
echo $! >> $pidfile

docker logs -f genesis_ganachecli_1 > ${datadir}/ganachelog.txt 2>&1 &
echo $! >> $pidfile

docker exec -ti integration_sifnode1_1 tail -f /tmp/testrun.sh > ${datadir}/clicmds.txt &
echo $! >> $pidfile
