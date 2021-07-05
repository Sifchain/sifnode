#!/usr/bin/env bash

killall sifnoded sifnodecli

#sifnodecli rest-server &
#sifnoded start

sifnoded start --home ~/.sifnode-2 --p2p.laddr 0.0.0.0:27656  --grpc.address 0.0.0.0:9091 --address tcp://0.0.0.0:27660 --rpc.laddr tcp://127.0.0.1:27658 >> abci_2.log 2>&1  &
sifnoded start --home ~/.sifnode-1 --p2p.laddr 0.0.0.0:27655  --grpc.address 0.0.0.0:9090 --address tcp://0.0.0.0:27659 --rpc.laddr tcp://127.0.0.1:27657 >> abci_1.log 2>&1  &
rm -rf ~/.ibc-setup/last-queried-heights.json
#Reset connections
ibc-setup ics20 --mnemonic "race draft rival universe maid cheese steel logic crowd fork comic easy truth drift tomorrow eye buddy head time cash swing swift midnight borrow"
ibc-relayer start -i -v --poll 10
#Created channel:
#  localnet-1: transfer/channel-0 (connection-0)
#  localnet-2: transfer/channel-0 (connection-0)

#sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd
#sifnoded q bank balances sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd --node tcp://127.0.0.1:27658
#sifnoded q bank balances sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd --node tcp://127.0.0.1:27657
#
#sifnoded tx ibc-transfer transfer transfer channel-1 sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd 100rowan --node tcp://127.0.0.1:27657 --chain-id=localnet-1 --from=sif --log_level=debug --gas-prices=0.5rowan --keyring-backend test  --home ~/.sifnode-1
#sifnoded tx ibc-transfer transfer transfer channel-0 sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd 100ibc/E0263CEED41F926DCE9A805F0358074873E478B515A94DF202E6B69E29DA6178 --node tcp://127.0.0.1:27658 --chain-id=localnet-2 --from=sif --log_level=debug --gas-prices=0.5rowan --keyring-backend test  --home ~/.sifnode-2


Test 1
- Stop one network
  Relayer continues to sync , reports err from one network
  Catches up when the network is started again ( syncs blocks)
  Transfer works fine after network is synced

Test 2
- Stop both networks
  Relayer stops
  Does not start syncing when first network is started
  Starts syncing when both networks are available
  Transfer works fine after network is synced

Test 3
- Stop relayer
  Can be started again it catches up. ( requires mnemonic)
  Transfer works fine after network is synced

Test 4
- Stop both networks and relayer
  Can be started again (requires mnemonic)
  Transfer works fine after network is synced

Test 5
- Stop relayer and clean last-queried-heights.json
  Can be started again (requires mnemonic)
  Transfer works fine after network is synced

Test 5
- Stop relayer and clean app.yaml (remove existing channels)
  Relayer cannot be started .
  Create new channel ( requires mnemonic)
  Fresh channel established ( Observation : it creates the next connection name , connection-1 instead of connection-0 , so older identifiers are cached somehwere )
  Transfer works fine after network is synced (need to add proper channel id)
  Different channel Id creates a different token denom (ibc/<different hash>)






