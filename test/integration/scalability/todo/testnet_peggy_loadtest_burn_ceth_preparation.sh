#!/bin/bash -x

SIFCHAIN_LOCK_ADDRESS=$(sifnoded keys show testnet-peggy-loadtest-burn-eth --keyring-backend test -a)
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