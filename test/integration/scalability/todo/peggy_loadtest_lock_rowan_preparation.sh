#!/bin/bash -x

SIF_ADDRESS=$(sifnoded keys show $1 --keyring-backend test -a)

yarn \
    -s \
    --cwd /sifnode/smart-contracts \
    integrationtest:sendLockTx \
    --sifchain_address ${SIF_ADDRESS} \
    --symbol eth \
    --ethereum_private_key_env_var ETH_PRIVATE_KEY \
    --json_path /sifnode/smart-contracts/deployments/${DEPLOYMENT_NAME} \
    --gas estimate \
    --ethereum_network ${ETH_NETWORK} \
    --bridgebank_address ${BRIDGEBANK_ADDRESS} \
    --ethereum_address ${ETH_ADDRESS} \
    --amount 1
    # --amount 50000000000000000