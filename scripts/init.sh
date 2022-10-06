#!/usr/bin/env bash

### chain init script for development purposes only ###

make clean install
rm -rf ~/.sifnoded
sifnoded init test --chain-id=localnet -o

echo "Generating deterministic account - sif"
echo "race draft rival universe maid cheese steel logic crowd fork comic easy truth drift tomorrow eye buddy head time cash swing swift midnight borrow" | sifnoded keys add sif --recover --keyring-backend=test

echo "Generating deterministic account - akasha"
echo "hand inmate canvas head lunar naive increase recycle dog ecology inhale december wide bubble hockey dice worth gravity ketchup feed balance parent secret orchard" | sifnoded keys add akasha --recover --keyring-backend=test

echo "Generating deterministic account - alice"
echo "crunch enable gauge equip sadness venture volcano capable boil pole lounge because service level giggle decide south deposit bike antique consider olympic girl butter" | sifnoded keys add alice --recover --keyring-backend=test

sifnoded keys add mkey --multisig sif,akasha --multisig-threshold 2 --keyring-backend=test

sifnoded add-genesis-account $(sifnoded keys show sif -a --keyring-backend=test) 500000000000000000000000000000000rowan,500000000000000000000000catk,500000000000000000000000cbtk,500000000000000000000000000000000ceth,990000000000000000000000000stake,500000000000000000000000cdash,500000000000000000000000clink,5000000000000cusdt,90000000000000000000ibc/96D7172B711F7F925DFC7579C6CCC3C80B762187215ABD082CDE99F81153DC80 --keyring-backend=test
sifnoded add-genesis-account $(sifnoded keys show akasha -a --keyring-backend=test) 500000000000000000000000rowan,500000000000000000000000catk,500000000000000000000000cbtk,500000000000000000000000ceth,990000000000000000000000000stake,500000000000000000000000cdash,500000000000000000000000clink --keyring-backend=test
sifnoded add-genesis-account $(sifnoded keys show alice -a --keyring-backend=test) 500000000000000000000000rowan,500000000000000000000000catk,500000000000000000000000cbtk,500000000000000000000000ceth,990000000000000000000000000stake,500000000000000000000000cdash,500000000000000000000000clink --keyring-backend=test

sifnoded add-genesis-clp-admin $(sifnoded keys show sif -a --keyring-backend=test) --keyring-backend=test
sifnoded add-genesis-clp-admin $(sifnoded keys show akasha -a --keyring-backend=test) --keyring-backend=test

sifnoded set-genesis-oracle-admin sif --keyring-backend=test
sifnoded add-genesis-validators $(sifnoded keys show sif -a --bech val --keyring-backend=test) --keyring-backend=test

sifnoded set-genesis-whitelister-admin sif --keyring-backend=test
sifnoded set-gen-denom-whitelist scripts/denoms.json

sifnoded gentx sif 1000000000000000000000000stake --moniker sif_val --chain-id=localnet --keyring-backend=test

echo "Collecting genesis txs..."
sifnoded collect-gentxs

echo "Validating genesis file..."
sifnoded validate-genesis
