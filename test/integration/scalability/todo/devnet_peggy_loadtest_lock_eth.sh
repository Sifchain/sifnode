#!/bin/bash -x

SIFCHAIN_LOCK_ADDRESS=$(sifnoded keys show devnet-peggy-loadtest-lock-eth --keyring-backend test -a)
for i in {1..100}
do
    yarn \
        -s \
        --cwd /sifnode/smart-contracts \
        integrationtest:sendLockTx \
        --sifchain_address ${SIFCHAIN_LOCK_ADDRESS} \
        --symbol eth \
        --ethereum_private_key_env_var ETHEREUM_PRIVATE_KEY \
        --json_path /sifnode/smart-contracts/deployments/sifchain-devnet-042 \
        --gas estimate \
        --ethereum_network ropsten \
        --bridgebank_address 0x471e0ffB16C4eEde754cEfD7F522257df37a1410 \
        --ethereum_address 0x5171050beb52148aB834Fb21E3E30FA429470c46 \
        --amount 1
done