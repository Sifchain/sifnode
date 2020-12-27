#!/bin/bash

basedir=$(dirname $0)
. $basedir/vagrantenv.sh

bash $basedir/sifchain_logs.sh
docker logs -t genesis_ganachecli_1 > ${datadir}/ganachelog.txt 2>&1
cp $SMART_CONTRACTS_DIR/.env ${datadir}/env
cp $TEST_INTEGRATION_DIR/vagrantenv.sh ${datadir}/vagrantenv.sh
( cd $SMART_CONTRACTS_DIR && truffle networks ) > ${datadir}/trufflenetworks.txt
sudo rsync -a $GANACHE_DB_DIR/ ${datadir}/ganachedb/ && chown -R $(id -u) ${datadir}/ganachedb/
