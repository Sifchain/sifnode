#!/usr/bin/env bash

killall sifnoded sifnodecli

sifnodecli rest-server &
sifnoded start

sifnoded start --home ~/.sifnode-2 --p2p.laddr 0.0.0.0:26656  --grpc.address 0.0.0.0:9091 --address tcp://0.0.0.0:26660 --rpc.laddr tcp://127.0.0.1:26658
sifnoded start --home ~/.sifnode-1 --p2p.laddr 0.0.0.0:26655  --grpc.address 0.0.0.0:9090 --address tcp://0.0.0.0:26659 --rpc.laddr tcp://127.0.0.1:26657