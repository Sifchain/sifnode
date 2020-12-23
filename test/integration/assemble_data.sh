#!/bin/bash

basedir=$(dirname $0)
datadir=${basedir}/vagrant/data

bash $basedir/sifchain_logs.sh
cp /sifnode/smart-contracts/.env ${datadir}/env
cp /sifnode/test/integration/vagrantenv.sh ${datadir}/vagrantenv.sh
touch /tmp/bridgebank.txt && cp /tmp/bridgebank.txt ${datadir}/bridgebank.txt
( cd /sifnode/smart-contracts && truffle networks ) > ${datadir}/trufflenetworks.txt
