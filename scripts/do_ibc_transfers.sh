#!/bin/zsh

sifnoded tx ibc-transfer transfer transfer channel-1 sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5 100rowan --node tcp://127.0.0.1:27666 --chain-id=localnet-2 --from=sif --log_level=debug  --keyring-backend test --fees 10000000000000000rowan  --home ~/.sifnode-2 --yes --broadcast-mode block
echo "Tried localnet-2 -> localnet-3"
echo ""

sifnoded tx ibc-transfer transfer transfer channel-0 sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5 100rowan --node tcp://127.0.0.1:27666 --chain-id=localnet-2 --from=sif --log_level=debug  --keyring-backend test --fees 10000000000000000rowan  --home ~/.sifnode-2 --yes --broadcast-mode block
echo "Tried localnet-2 -> localnet-1"
echo ""

sifnoded tx ibc-transfer transfer transfer channel-0 sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5 100rowan --node tcp://127.0.0.1:27665 --chain-id=localnet-1 --from=sif --log_level=debug  --keyring-backend test --fees 10000000000000000rowan  --home ~/.sifnode-1 --yes --broadcast-mode block
echo "Tried localnet-1 -> localnet-2"
echo ""

sifnoded tx ibc-transfer transfer transfer channel-1 sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5 100rowan --node tcp://127.0.0.1:27667 --chain-id=localnet-3 --from=sif --log_level=debug  --keyring-backend test --fees 10000000000000000rowan  --home ~/.sifnode-3 --yes --broadcast-mode block
echo "Tried localnet-3 -> localnet-1"
echo ""

sifnoded tx ibc-transfer transfer transfer channel-1 sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5 100rowan --node tcp://127.0.0.1:27665 --chain-id=localnet-1 --from=sif --log_level=debug  --keyring-backend test --fees 10000000000000000rowan  --home ~/.sifnode-1 --yes --broadcast-mode block
echo "Tried localnet-1 -> localnet-3"
