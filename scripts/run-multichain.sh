#!/usr/bin/env bash

killall sifnoded sifnodecli

#sifnodecli rest-server &
#sifnoded start

sifnoded start --home ~/.sifnode-2 --p2p.laddr 0.0.0.0:27656  --grpc.address 0.0.0.0:9091 --address tcp://0.0.0.0:27660 --rpc.laddr tcp://127.0.0.1:27658 >> abci_2.log 2>&1  &
sifnoded start --home ~/.sifnode-1 --p2p.laddr 0.0.0.0:27655  --grpc.address 0.0.0.0:9090 --address tcp://0.0.0.0:27659 --rpc.laddr tcp://127.0.0.1:27657 >> abci_1.log 2>&1  &
rm -rf ~/.ibc-setup/last-queried-heights.json
#Reset connections
ibc-setup ics20 --mnemonic "race draft rival universe maid cheese steel logic crowd fork comic easy truth drift tomorrow eye buddy head time cash swing swift midnight borrow" --home ~/.ibc-12
ibc-setup ics20 --mnemonic "race draft rival universe maid cheese steel logic crowd fork comic easy truth drift tomorrow eye buddy head time cash swing swift midnight borrow" --home ~/.ibc-31
ibc-relayer start -i -v --poll 10 --home ~/.ibc-12
ibc-relayer start -i -v --poll 10 --home ~/.ibc-23
ibc-relayer start -i -v --poll 10 --home ~/.ibc-31
#Created channel:
#  localnet-1: transfer/channel-0 (connection-0)
#  localnet-2: transfer/channel-0 (connection-0)

#sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd
#sifnoded q bank balances sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd --node tcp://127.0.0.1:27665
#sifnoded q bank balances sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5 --node tcp://127.0.0.1:27666
#sifnoded q bank balances sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5 --node tcp://127.0.0.1:27667
#
#sifnoded tx ibc-transfer transfer transfer channel-0 sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5 100rowan --node tcp://127.0.0.1:27665 --chain-id=localnet-1 --from=akasha --log_level=debug --gas-prices=0.5rowan --keyring-backend test  --home ~/.sifnode-1
#sifnoded tx ibc-transfer transfer transfer channel-2 sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5 50ibc/E0263CEED41F926DCE9A805F0358074873E478B515A94DF202E6B69E29DA6178 --node tcp://127.0.0.1:27666 --chain-id=localnet-2 --from=akasha --log_level=debug --gas-prices=0.5rowan --keyring-backend test  --home ~/.sifnode-2
#sifnoded tx ibc-transfer transfer transfer channel-0 sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd 50ibc/4C2B3D3B398FC7B8FFA3A96314006FF0B38E3BFC4CE90D8EE46E9EB6768A482D --node tcp://127.0.0.1:27666 --chain-id=localnet-2 --from=sif --log_level=debug --gas-prices=0.5rowan --keyring-backend test  --home ~/.sifnode-2
#sifnoded tx ibc-transfer transfer transfer channel-1 sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5 50ibc/5C3977A32007D22B1845B57076D0E27C3159C3067B11B9CEF6FA551D71DAEDD6 --node tcp://127.0.0.1:27667 --chain-id=localnet-3 --from=akasha --log_level=debug --gas-prices=0.5rowan --keyring-backend test  --home ~/.sifnode-3


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






Transfer 1 [1-2]
Source Port    :  transfer
Source Channel :  channel-1
rpc error: code = InvalidArgument desc = invalid denom trace hash rowan, encoding/hex: invalid byte: U+0072 'r'
Trace : <nil>
encoding/hex: invalid byte: U+006E 'n': invalid denomination for cross-chain transfer [cosmos/cosmos-sdk@v0.42.4/x/ibc/applications/transfer/keeper/relay.go:396]
PATH :
Denom :  rowan

Transfer 2 [2->3]
Source Port    :  transfer
Source Channel :  channel-0
rpc error: code = InvalidArgument desc = invalid denom trace hash transfer/channel-1/rowan, encoding/hex: invalid byte: U+0074 't'
Trace : <nil>
encoding/hex: invalid byte: U+0073 's': invalid denomination for cross-chain transfer [cosmos/cosmos-sdk@v0.42.4/x/ibc/applications/transfer/keeper/relay.go:396]
PATH :
Denom :  transfer/channel-1/rowan  [channel-1 is the channel used in trasfer from 1->2]

Trasfer 3[3-1]
Source Port    :  transfer
Source Channel :  channel-1
rpc error: code = InvalidArgument desc = invalid denom trace hash transfer/channel-0/transfer/channel-1/rowan, encoding/hex: invalid byte: U+0074 't'
Trace : <nil>
encoding/hex: invalid byte: U+0073 's': invalid denomination for cross-chain transfer [cosmos/cosmos-sdk@v0.42.4/x/ibc/applications/transfer/keeper/relay.go:396]
PATH :
Denom :  transfer/channel-0/transfer/channel-1/rowan


Trasfer 4 -[2-1]
Source Port    :  transfer
Source Channel :  channel-1
rpc error: code = InvalidArgument desc = invalid denom trace hash transfer/channel-1/rowan, encoding/hex: invalid byte: U+0074 't'
Trace : <nil>
encoding/hex: invalid byte: U+0073 's': invalid denomination for cross-chain transfer [cosmos/cosmos-sdk@v0.42.4/x/ibc/applications/transfer/keeper/relay.go:396]
PATH :
Denom :  transfer/channel-1/rowan





Transfer  1-2 : rowan : 4CE5019F09EDEBC793F1B90547638F6E00B5E5F4CB37D10585EAF9A3981376FF  / minted : E0263CEED41F926DCE9A805F0358074873E478B515A94DF202E6B69E29DA6178
Transfer  2-3 : transfer/channel-0/rowan : E0263CEED41F926DCE9A805F0358074873E478B515A94DF202E6B69E29DA6178 minted : 5C3977A32007D22B1845B57076D0E27C3159C3067B11B9CEF6FA551D71DAEDD6
Transfer  2-1 : Returning to source | Denom : transfer/channel-0/rowan , SourcePort : transfer , SourceChannel : channel-0
TRasfer   3-1 : transfer/channel-0/transfer/channel-0/rowan : 5C3977A32007D22B1845B57076D0E27C3159C3067B11B9CEF6FA551D71DAEDD6 minted : F3A091ED158AE3605D066E6A27AAEE9AD98ABEA8894CF9F2621D20A3CD381BE3
