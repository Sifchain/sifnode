#!/bin/zsh

echo "initializing sifnoded and doing initial IBC transfers"
scripts/init-multichain.sh > /dev/null 2>&1
scripts/do_ibc_transfers.sh > /dev/null 2>&1

echo "killing sifnode-2"
ps auxww | grep sifnode-2 | egrep -v grep | awk '{ print $2 }' | xargs kill

sleep 60
echo "killing hermes"
killall hermes

# start hermes
echo "starting hermes"
hermes start > hermes.log 2>&1 &

echo "Sleeping to let hermes boot"
sleep 10

echo "second round of ibc transfers...some should fail so ignore errors"
scripts/do_ibc_transfers.sh

echo "sleeping to let the localnet-2 expire"
sleep 850

echo "restarting localnet-2 (sifnode)"
sifnoded start --home ~/.sifnode-2 --p2p.laddr 0.0.0.0:27656  --grpc.address 0.0.0.0:9091 --grpc-web.address 0.0.0.0:9094 --address tcp://0.0.0.0:27660 --rpc.laddr tcp://127.0.0.1:27666 >> abci_2.log 2>&1  &
sleep 10

echo "bouncing hermes"
killall hermes

# start hermes
echo "starting hermes"
hermes start > hermes.log 2>&1 &

echo "Sleeping to let hermes boot"
sleep 10

echo "Checking channels"
hermes query packet unreceived-packets localnet-1 transfer channel-0
hermes query packet unreceived-packets localnet-1 transfer channel-1
hermes query packet unreceived-packets localnet-2 transfer channel-0
hermes query packet unreceived-packets localnet-2 transfer channel-1
hermes query packet unreceived-packets localnet-3 transfer channel-0
hermes query packet unreceived-packets localnet-3 transfer channel-1

echo "Creating channels to clone"
scripts/create_clonable_hermes_channels.sh > /dev/null 2>&1

echo "Current Terra and balances"
echo "localnet-1"
terrad q bank balances $(terrad keys show terra -a --keyring-backend=test --home ~/.terranode-1) --node tcp://127.0.0.1:27665
echo ""
echo "localnet-2"
sifnoded q bank balances $(sifnoded keys show sif -a --keyring-backend=test --home ~/.sifnode-2) --node tcp://127.0.0.1:27666
echo ""
echo "localnet-3"
sifnoded q bank balances $(sifnoded keys show sif -a --keyring-backend=test --home ~/.sifnode-2) --node tcp://127.0.0.1:27667
echo ""

echo "Current Akasha balances"
echo "localnet-1"
terrad q bank balances $(terrad keys show akasha -a --keyring-backend=test --home ~/.terranode-1) --node tcp://127.0.0.1:27665
echo ""
echo "localnet-2"
sifnoded q bank balances $(sifnoded keys show akasha -a --keyring-backend=test --home ~/.sifnode-2) --node tcp://127.0.0.1:27666
echo ""
echo "localnet-3"
sifnoded q bank balances $(sifnoded keys show akasha -a --keyring-backend=test --home ~/.sifnode-2) --node tcp://127.0.0.1:27667
echo ""

echo "voting on proposals"
#scripts/vote.sh > /dev/null 2>&1
#scripts/vote_localnet_1.sh > /dev/null 2>&1
