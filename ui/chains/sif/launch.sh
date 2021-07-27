#!/usr/bin/env bash

. ../credentials.sh

rm -rf ~/.sifnoded

sifnoded init test --chain-id=sifchain-local
cp ./app.toml ~/.sifnoded/config

#sifnoded config output json
#sifnoded config indent true
#sifnoded config trust-node true
#sifnoded config chain-id sifchain-local
#sifnoded config keyring-backend test

echo "Generating deterministic account - ${SHADOWFIEND_NAME}"
echo "${SHADOWFIEND_MNEMONIC}" | sifnoded keys add ${SHADOWFIEND_NAME}  --keyring-backend=test --recover

echo "Generating deterministic account - ${AKASHA_NAME}"
echo "${AKASHA_MNEMONIC}" | sifnoded keys add ${AKASHA_NAME}  --keyring-backend=test --recover

echo "Generating deterministic account - ${JUNIPER_NAME}"
echo "${JUNIPER_MNEMONIC}" | sifnoded keys add ${JUNIPER_NAME} --keyring-backend=test --recover

sifnoded add-genesis-account $(sifnoded keys show ${SHADOWFIEND_NAME} -a --keyring-backend=test) 100000000000000000000000000000rowan,100000000000000000000000000000catk,100000000000000000000000000000cbtk,100000000000000000000000000000ceth,100000000000000000000000000000cusdc,100000000000000000000000000000clink,100000000000000000000000000stake
sifnoded add-genesis-account $(sifnoded keys show ${AKASHA_NAME} -a --keyring-backend=test) 100000000000000000000000000000rowan,100000000000000000000000000000catk,100000000000000000000000000000cbtk,100000000000000000000000000000ceth,100000000000000000000000000000cusdc,100000000000000000000000000000clink,100000000000000000000000000stake
sifnoded add-genesis-account $(sifnoded keys show ${JUNIPER_NAME} -a --keyring-backend=test) 10000000000000000000000rowan,10000000000000000000000cusdc,100000000000000000000clink,100000000000000000000ceth

sifnoded add-genesis-validators $(sifnoded keys show ${SHADOWFIEND_NAME} -a --bech val --keyring-backend=test)

sifnoded gentx ${SHADOWFIEND_NAME} 1000000000000000000000000stake --chain-id=sifchain-local --keyring-backend test

echo "Collecting genesis txs..."
sifnoded collect-gentxs

echo "Validating genesis file..."
sifnoded validate-genesis

echo "Starting test chain"

./start.sh
