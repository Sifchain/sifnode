#!/bin/bash -x

ACCOUNT_NAME="testnet-peggy-loadtest-lock-rowan-2021-08-16"
SIFCHAIN_LOCK_ADDRESS=$(sifnoded keys show testnet-peggy-loadtest-lock-rowan-2021-08-16 --keyring-backend test -a)
yarn \
    -s \
    --cwd /sifnode/smart-contracts \
    integrationtest:sendLockTx \
    --sifchain_address ${SIFCHAIN_LOCK_ADDRESS} \
    --symbol eth \
    --ethereum_private_key_env_var ETHEREUM_PRIVATE_KEY \
    --json_path /sifnode/smart-contracts/deployments/sifchain-testnet-042-ibc \
    --gas estimate \
    --ethereum_network ropsten \
    --bridgebank_address 0xB75849afEF2864977a858073458Cb13F9410f8e5 \
    --ethereum_address 0x5171050beb52148aB834Fb21E3E30FA429470c46 \
    --amount 4000000000000100000
./testnet_send.sh testnet-source testnet-peggy-loadtest-lock-rowan 10000000000000rowan