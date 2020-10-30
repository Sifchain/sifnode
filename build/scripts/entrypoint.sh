#!/bin/sh

sifnodecli rest-server --laddr tcp://0.0.0.0:1317 &
sifnoded start --rpc.laddr tcp://0.0.0.0:26657
# ebrelayer init tcp://0.0.0.0:26656  ws://localhost:7545/ $($PEGGY_CONTRACT_ADDRESS) $($MONIKER) --chain-id=$(CHAINNET)
